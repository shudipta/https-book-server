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
)

const (
	GettingSecretError      = 1
	DecodingConfigDataError = 2
)

type RegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetAllSecrets(refs []corev1.LocalObjectReference) []string {
	var names []string
	for _, ref := range refs {
		names = append(names, ref.Name)
	}
	return names
}

func (c *ScannerController) CheckWorkload(w *workload.Workload) (*workload.Workload, bool, error) {
	secretNames := GetAllSecrets(w.Spec.Template.Spec.ImagePullSecrets)

	allow, err := c.CheckContainers(w.ObjectMeta.GetNamespace(), w.Spec.Template.Spec.InitContainers, secretNames)
	if !allow {
		return w, false, err
	}

	allow, err = c.CheckContainers(w.ObjectMeta.GetNamespace(), w.Spec.Template.Spec.Containers, secretNames)
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
	namespace string, containers []corev1.Container, secretNames []string) (bool, error) {
	for _, cont := range containers {
		_, _, status, err := c.CheckImage(namespace, cont.Image, secretNames)
		vulnerable := status != scanner.NotVulnerableStatus
		if vulnerable {
			return false, err
		}
	}
	return true, nil
}

// This method takes namespace_name <namespace> of provided secrets <secretNames> and image name
// of a docker image. For each secret, it reads the config data of secret and store it to
// auth variable (map[string]map[string]RegistrySecret)
// we need this type to store config data, because original config date is in following format:
// {
//   "auths":{
// 	   <api url>:{
// 	 	 "username":<username>,
// 	 	 "password":<password>,
// 	 	 "email":<email>,
// 	 	 "auth":<auth token>
// 	   }
// 	 }
// }
// Then it scans to find vulnerabilities in the image for all credentials. It returns
// 			(true, error); if any error occured
// 			(false, nil); if no vulnerability exists
// If the image is not found with the secret info, then it tries with the public docker
// url="https://registry-1.docker.io/"
func (c *ScannerController) CheckImage(
	namespace, image string,
	secretNames []string) ([]api.Feature, []api.Vulnerability, int, error) {

	for _, item := range secretNames {
		secret, err := c.KubeClient.CoreV1().Secrets(namespace).Get(item, metav1.GetOptions{})
		if err != nil {
			return []api.Feature{}, []api.Vulnerability{}, GettingSecretError,
				fmt.Errorf("error in reading secret(%s): \n\t%v", item, err)
		}

		var configData []byte
		for _, val := range secret.Data {
			configData = append(configData, val...)
			break
		}

		var authInfo map[string]map[string]RegistrySecret
		err = json.NewDecoder(bytes.NewReader(configData)).Decode(&authInfo)
		if err != nil {
			return []api.Feature{}, []api.Vulnerability{}, DecodingConfigDataError,
				fmt.Errorf("error in decoding configData of secret(%s): \n\t%v", item, err)
		}

		for _, authInfo := range authInfo {
			for key, val := range authInfo {
				features, vulnerabilities, status, err := scanner.IsVulnerable(
					c.KubeClient,
					key, image, val.Username, val.Password,
				)
				if status > 4 {
					return features, vulnerabilities, status, err
				}
			}
		}
	}

	registryUrl := "https://registry-1.docker.io/"
	username := "" // anonymous
	password := "" // anonymous

	features, vulnerabilities, status, err := scanner.IsVulnerable(
		c.KubeClient,
		registryUrl, image, username, password,
	)
	if status < 5 {
		return features, vulnerabilities, status, fmt.Errorf("error in secrets for image(%s): %v", image, err)
	}

	return features, vulnerabilities, status, err
}
