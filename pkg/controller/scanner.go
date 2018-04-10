package controller

import (
	"bytes"
	"encoding/json"
	"fmt"

	workload "github.com/appscode/kubernetes-webhook-util/workload/v1"
	"github.com/soter/scanner/pkg/scanner"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/util/parsers"
)

//type AuthInfo struct {
//	Auth map[string]RegistrySecret `json:"auth"`
//}

type RegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// getAllSecrets() takes imagePullSecrets and return the list of secret names as an array of
// string
func getAllSecrets(imagePullSecrets []corev1.LocalObjectReference) []string {
	secretNames := []string{}
	for _, secretName := range imagePullSecrets {
		secretNames = append(secretNames, secretName.Name)
	}

	return secretNames
}

// checkworkload() checks vulnerabilities in the images used in containes
func (c *ScannerController) checkWorkload(w *workload.Workload) (*workload.Workload, bool, error) {
	secretNames := getAllSecrets(w.Spec.Template.Spec.ImagePullSecrets)

	for _, cont := range w.Spec.Template.Spec.InitContainers {
		vulnerable, err := c.checkImage(w.ObjectMeta.GetNamespace(), cont.Image, secretNames)
		if vulnerable {
			return w, false, err
		}

	}
	for _, cont := range w.Spec.Template.Spec.Containers {
		vulnerable, err := c.checkImage(w.ObjectMeta.GetNamespace(), cont.Image, secretNames)
		if vulnerable {
			return w, false, err
		}
	}

	return w, true, nil
}

// This method takes namespace_name <namespace> of provided secrets <secretNames> and image name
// of a docker image. For each secret, it reads the config data of secret and store it to
// registrySecrets (map[string]RegistrySecret) where the api url is the key and value is the
// credentials. Then it scans to find vulnerabilities in the image for all secrets' content. It returns
// 			(true, error); if any error occured
// 			(false, nil); if no vulnerability exists
// If the image is not found with the secret info, then it tries with the public docker
// url="https://registry-1.docker.io/"
func (c *ScannerController) checkImage(namespace, image string, secretNames []string) (bool, error) {
	// TODO: check for digest
	repoName, tag, _, err := parsers.ParseImageName(image)
	if err != nil {
		return true, fmt.Errorf("invalid image(%s), %v", image, err)
	}

	// Here repo has "docker.io/" as a prefix. So we drop it
	repoName = repoName[10:]

	for _, item := range secretNames {
		secret, err := c.KubeClient.CoreV1().Secrets(namespace).Get(item, metav1.GetOptions{})
		if err != nil {
			return true, fmt.Errorf("error in reading secret(%s): \n\t%v", item, err)
		}

		configData := []byte{}
		for _, val := range secret.Data {
			configData = append(configData, val...)
			break
		}

		var auth map[string]map[string]RegistrySecret
		err = json.NewDecoder(bytes.NewReader(configData)).Decode(&auth)
		if err != nil {
			return true, fmt.Errorf("error in decoding configData of secret(%s): \n\t%v", item, err)
		}

		for _, authInfo := range auth {
			for key, val := range authInfo {
				vulnerable, status, err := scanner.IsVulnerable(c.KubeClient, key, repoName, tag, val.Username, val.Password)
				if status > 2 {
					return vulnerable, err
				}
			}
		}
	}

	registryUrl := "https://registry-1.docker.io/"
	username := "" // anonymous
	password := "" // anonymous

	vulnerable, status, err := scanner.IsVulnerable(c.KubeClient, registryUrl, repoName, tag, username, password)
	if status < 3 {
		return true, fmt.Errorf("error in secrets for image(%s): %v", image, err)
	}

	return vulnerable, err
}
