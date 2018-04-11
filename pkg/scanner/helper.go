package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/golang-lru"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// requestBearerToken() method create a http.request object from docker for given credential
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

// getBearerToken() takes a http.Response and makes the bearer token from response body
// by adding "Bearer " as prefix to it
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

// getAddress() forms the clairAddress to call rest clair api
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

// requestSendingLayer() takes layer object <l> and <clairAddr> and creates a http.Request
// to send this layer to clair running at <clairAddr>
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

// requestVulnerabilities() creates a http.Request for an image indicated by it's last
// layer's hash(digest) to get the vulnerabilities. This request is sent to <clairAddr>
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

// decode() method just decode the response body into a layerApi object
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

// getVulnerabilities() collects vulnerabilities if exist in the layer
func getVulnerabilities(layerObj layerApi) []Vulnerability {
	vuls := []Vulnerability{}
	for _, feature := range layerObj.Layer.Features {
		for _, vul := range feature.Vulnerabilities {
			vuls = append(vuls, vul)
		}
	}

	return vuls
}

// getFeatures() collects Features in the layer
func getFeatures(layerObj layerApi) []Feature {
	fs := []Feature{}
	for _, feature := range layerObj.Layer.Features {
		fs = append(fs, Feature{feature.Name, feature.NamespaceName, feature.Version})
	}

	return fs
}

// parseImageName() process a full qualified image name we used in containers
// and returns three parts(registry, repository, tag) for this image
func parseImageName(image string) (string, string, string) {
	registry := "registry-1.docker.io"
	tag := "latest"
	var nameParts, tagParts []string
	var name, port string
	state := 0
	start := 0
	for i, c := range image {
		if c == ':' || c == '/' || c == '@' || i == len(image)-1 {
			if i == len(image)-1 {
				i += 1
			}
			part := image[start:i]
			start = i + 1
			switch state {
			case 0:
				if strings.Contains(part, ".") {
					registry = part
					if c == ':' {
						state = 1
					} else {
						state = 2
					}
				} else {
					if c == '/' {
						start = 0
						state = 2
					} else {
						state = 3
						name = fmt.Sprintf("library/%s", part)
					}
				}
			case 3:
				tag = ""
				tagParts = append(tagParts, part)
			case 1:
				state = 2
				port = part
			case 2:
				if c == ':' || c == '@' {
					state = 3
				}
				nameParts = append(nameParts, part)
			}
		}
	}

	if port != "" {
		registry = fmt.Sprintf("%s:%s", registry, port)
	}

	if name == "" {
		name = strings.Join(nameParts, "/")
	}

	if tag == "" {
		tag = strings.Join(tagParts, ":")
	}

	registry = fmt.Sprintf("https://%s", registry)

	return registry, name, tag
}

func HashPart(digest string) string {
	if len(digest) < 7 {
		return ""
	}

	return digest[7:]
}

// cacheFeaturesAndVulnerabilities() just cache the given vulnerabilities
func cacheFeaturesAndVulnerabilities(fsCache, vulsCache *lru.TwoQueueCache, layerName string, fs []Feature, vuls []Vulnerability) {
	fsCache.Add(layerName, fs)
	vulsCache.Add(layerName, vuls)
}
