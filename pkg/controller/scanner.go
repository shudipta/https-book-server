package controller

import (
	"bytes"
	"encoding/json"
	"fmt"

	workload "github.com/appscode/kubernetes-webhook-util/workload/v1"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"github.com/soter/scanner/pkg/scanner"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	GettingSecretError      = 1
	DecodingConfigDataError = 2
)

type RegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// getAllSecrets() takes imagePullSecrets and return the list of secret names as an array of
// string
func GetAllSecrets(imagePullSecrets []corev1.LocalObjectReference) []string {
	secretNames := []string{}
	if (imagePullSecrets) == nil {
		return []string{}
	}
	for _, secretName := range imagePullSecrets {
		secretNames = append(secretNames, secretName.Name)
	}

	return secretNames
}

// checkworkload() checks vulnerabilities for given workload obj
func (c *ScannerController) CheckWorkload(w *workload.Workload, precache bool) (*workload.Workload, bool, error) {
	secretNames := GetAllSecrets(w.Spec.Template.Spec.ImagePullSecrets)

	allow, err := c.CheckContainers(w.ObjectMeta.GetNamespace(), w.Spec.Template.Spec.InitContainers, secretNames, precache)
	if !allow {
		return w, false, err
	}

	allow, err = c.CheckContainers(w.ObjectMeta.GetNamespace(), w.Spec.Template.Spec.Containers, secretNames, precache)
	if !allow {
		return w, false, err
	}

	return w, true, nil
}

// checkContainers() checks vulnerabilities for each images used in containers.
// Here, precache parameter indicates that checking is being done for storing
// vulnerabilities and features of each image layer into cache. Otherwise,
// if precache is false then
// 		if any image is vulnerable then
//           this method returns
func (c *ScannerController) CheckContainers(
	namespace string, containers []corev1.Container, secretNames []string,
	precache bool) (bool, error) {
	for _, cont := range containers {
		_, status, err := c.CheckImage(namespace, cont.Image, secretNames, precache)
		vulnerable := (status != scanner.NotVulnerableStatus)
		if !precache && vulnerable {
			return false, err
		}
	}

	return true, nil
}

// This method takes namespace_name <namespace> of provided secrets <secretNames> and image name
// of a docker image. For each secret, it reads the config data of secret and store it to
// registrySecrets (map[string]RegistrySecret) where the api url is the key and value is the
// credentials. Then it scans to find vulnerabilities in the image for all secrets' content. It returns
// 			(true, error); if any error occured
// 			(false, nil); if no vulnerability exists
// If the image is not found with the secret info, then it tries with the public docker
// url="https://registry-1.docker.io/"
func (c *ScannerController) CheckImage(
	namespace, image string,
	secretNames []string, precache bool) (scanner.Canonical1, int, error) {

	if (secretNames) == nil {
		return scanner.Canonical1{}, 0, fmt.Errorf("empty SecretNames[]")
	}
	for _, item := range secretNames {
		secret, err := c.KubeClient.CoreV1().Secrets(namespace).Get(item, metav1.GetOptions{})
		if err != nil {
			return scanner.Canonical1{}, GettingSecretError,
				fmt.Errorf("error in reading secret(%s): \n\t%v", item, err)
		}

		configData := []byte{}
		for _, val := range secret.Data {
			configData = append(configData, val...)
			break
		}

		var authInfo map[string]map[string]RegistrySecret
		err = json.NewDecoder(bytes.NewReader(configData)).Decode(&authInfo)
		if err != nil {
			return scanner.Canonical1{}, DecodingConfigDataError,
				fmt.Errorf("error in decoding configData of secret(%s): \n\t%v", item, err)
		}

		for _, authInfo := range authInfo {
			for key, val := range authInfo {
				imageManifest, status, err := scanner.IsVulnerable(
					c.KubeClient, c.FsCache, c.VulsCache,
					key, image, val.Username, val.Password,
					precache,
				)
				if status > 3 {
					return imageManifest, status, err
				}
			}
		}
	}

	//registryUrl := "https://registry-1.docker.io/"
	registryUrl := ""
	username := "" // anonymous
	password := "" // anonymous

	imageManifest, status, err := scanner.IsVulnerable(
		c.KubeClient, c.FsCache, c.VulsCache,
		registryUrl, image, username, password,
		precache,
	)
	if status < 4 {
		return imageManifest, status, fmt.Errorf("error in secrets for image(%s): %v", image, err)
	}

	return imageManifest, status, err
}

type Feature map[api.Feature]struct{}
type Vulnerability map[api.Vulnerability]struct{}

func (c *ScannerController) GetImageReview(imageManifest scanner.Canonical1) ([]api.Feature, []api.Vulnerability) {
	var (
		fs        []api.Feature
		fsNameSet = sets.NewString()
		fsSet     = Feature{}

		vuls        []api.Vulnerability
		vulsNameSet = sets.NewString()
		vulsSet     = Vulnerability{}
	)

	for _, layer := range imageManifest.Layers {
		key := scanner.HashPart(imageManifest.Config.Digest) + scanner.HashPart(layer.Digest)

		valF, _ := c.FsCache.Get(key)
		fs1 := valF.([]api.Feature)
		if len(fs1) > 0 {
			for _, f := range fs1 {
				name := f.Name + f.NamespaceName + f.Version
				if !fsNameSet.Has(name) {
					fsNameSet.Insert(name)
					fsSet[f] = struct{}{}
				}
			}
		}

		valV, _ := c.VulsCache.Get(key)
		vuls1 := valV.([]api.Vulnerability)
		if len(vuls1) > 0 {
			for _, v := range vuls1 {
				name := v.Name + v.NamespaceName
				if !vulsNameSet.Has(name) {
					vulsNameSet.Insert(name)
					vulsSet[v] = struct{}{}
				}
			}
		}
	}

	for f := range fsSet {
		fs = append(fs, f)
	}
	for v := range vulsSet {
		vuls = append(vuls, v)
	}

	return fs, vuls
}
