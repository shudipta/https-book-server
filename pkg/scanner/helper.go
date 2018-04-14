package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func requestBearerToken(repo, userName, password string) (*http.Request, error) {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repo + ":pull&account=" + userName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if userName != "" {
		req.SetBasicAuth(userName, password)
	}

	return req, nil
}

func getBearerToken(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var token struct {
		Token string
	}

	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", token.Token), nil
}

func getAddress(kubeClient kubernetes.Interface) (string, error) {
	var host, port string
	pods, err := kubeClient.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, "clair") {
			host = pod.Status.HostIP
		}
	}

	clairSvc, err := kubeClient.CoreV1().Services("default").Get("clairsvc", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, p := range clairSvc.Spec.Ports {
		if p.TargetPort.IntVal == 6060 {
			port = strconv.Itoa(int(p.NodePort))
			break
		}
	}

	if host != "" && port != "" {
		return "http://" + host + ":" + port, nil
	}

	return "", fmt.Errorf("clair isn't running")
}

func requestSendingLayer(l *LayerType, clairAddr string) (*http.Request, error) {
	var layerApi struct {
		Layer *LayerType
	}
	layerApi.Layer = l
	reqBody, err := json.Marshal(layerApi)
	if err != nil {
		return nil, err
	}
	url := clairAddr + "/v1/layers"

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func requestVulnerabilities(hashNameOfImage, clairAddr string) (*http.Request, error) {
	url := clairAddr + "/v1/layers/" + hashNameOfImage + "?vulnerabilities"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

type layerApi struct {
	Layer *LayerType
}

func decode(resp *http.Response, err error) (layerApi, error) {
	if err != nil {
		return layerApi{}, err
	}
	defer resp.Body.Close()

	var layerObj layerApi
	err = json.NewDecoder(resp.Body).Decode(&layerObj)
	if err != nil {
		return layerApi{}, err
	} else if layerObj.Layer == nil {
		return layerApi{}, fmt.Errorf("clair returned empty layerObj")
	}

	return layerObj, nil
}

func getVulnerabilities(layerObj layerApi) []api.Vulnerability {
	var vuls []api.Vulnerability
	for _, feature := range layerObj.Layer.Features {
		for _, vul := range feature.Vulnerabilities {
			vuls = append(vuls, vul)
		}
	}

	return vuls
}

func getFeatures(layerObj layerApi) []api.Feature {
	var fs []api.Feature
	for _, feature := range layerObj.Layer.Features {
		fs = append(fs, api.Feature{feature.Name, feature.NamespaceName, feature.Version})
	}

	return fs
}

func parseImageName(imageName, registryUrl string) (string, string, string, string, error) {
	repo, tag, digest, err := parsers.ParseImageName(imageName)
	if err != nil {
		return "", "", "", "", err
	}
	// the repo part should have registry url as prefix followed by a '/'
	// for example, if image name = "ubuntu" then
	//					repo = "docker.io/library/ubuntu", tag = "latest", digest = ""
	// 				if image name = "k8s.gcr.io/kubernetes-dashboard-amd64:v1.8.1" then
	//					repo = "k8s.gcr.io/kubernetes-dashboard-amd64", tag = "v1.8.1", digest = ""
	// here, for docker registry the api url is "https://registry-1.docker.io"
	// and for other registry the url is "https://k8s.gcr.io"(gcr) or "https://quay.io"(quay)
	parts := strings.Split(repo, "/")
	if registryUrl == "" {
		if parts[0] == "docker.io" {
			registryUrl = "https://registry-1." + parts[0]
		} else {
			registryUrl = "https://" + parts[0]
		}
	}
	repo = strings.Join(parts[1:], "/")

	return registryUrl, repo, tag, digest, err
}

func HashPart(digest string) string {
	if len(digest) < 7 {
		return ""
	}

	return digest[7:]
}
