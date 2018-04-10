package scanner

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	reg "github.com/heroku/docker-registry-client/registry"
	oneliners "github.com/tamalsaha/go-oneliners"
	"k8s.io/client-go/kubernetes"
)

type config struct {
	MediaType string
	Size      int
	Digest    string
}

type layer struct {
	MediaType string
	Size      int
	Digest    string
}

type Canonical struct {
	SchemaVersion int
	MediaType     string
	Config        config
	Layers        []layer
}

const (
	RegistryConnectionError      = 1
	GettingManifestError         = 2
	GettingCannonicalError       = 3
	DecodingCannonicalError      = 4
	PullingLayersError           = 5
	BearerTokenRequestError      = 6
	BearerTokenResponseError     = 7
	GettingClairAddrError        = 8
	SendingLayerRequestError     = 9
	SendingLayerError            = 10
	VulnerabilitiesRequestError  = 11
	VulnerabilitiesResponseError = 12
	VulnerableStatus             = 13
	NotVulnerableStatus          = 14
)

func IsVulnerable(kubeClient kubernetes.Interface, registryUrl, repo, tag, username, password string) (bool, int, error) {
	imageName := repo + ":" + tag
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, //true
			},
		},
		Timeout: time.Minute,
	}

	hub, err := reg.New(registryUrl, username, password)
	if err != nil {
		return true, RegistryConnectionError,
			fmt.Errorf("error in connectiong to registry: %v\n", err)
	}

	manifest, err := hub.ManifestV2(repo, tag)
	if err != nil {
		return true, GettingManifestError,
			fmt.Errorf("error in getting manifest for image(%s): %v\n", imageName, err)
	}
	canonical, err := manifest.MarshalJSON()
	if err != nil {
		return true, GettingCannonicalError,
			fmt.Errorf("error in getting manifest.canonical for image(%s): %v\n", imageName, err)
	}
	canonicalReader := bytes.NewReader(canonical)

	var image Canonical
	if err := json.NewDecoder(canonicalReader).Decode(&image); err != nil {
		return true, DecodingCannonicalError,
			fmt.Errorf("error in decoding canonical for image(%s): %v\n", imageName, err)
	}
	oneliners.PrettyJson(image, "image")

	var layers []layer
	for _, l := range image.Layers {
		if l.Digest == "" {
			continue
		}
		layers = append(layers, l)
	}
	digest := image.Config.Digest

	layersLen := len(layers)
	if layersLen == 0 {
		return true, PullingLayersError,
			fmt.Errorf("error is pulling fsLayers for image(%s)\n", imageName)
	} else {
		fmt.Println("Analysing", len(layers), "layers")
	}

	clairClient := http.Client{
		Timeout: time.Minute,
	}

	var parent string
	req, _ := requestBearerToken(repo, username, password)
	if err != nil {
		return true, BearerTokenRequestError,
			fmt.Errorf("error in creating BearerToken request for image(%s): %v\n", imageName, err)
	}

	token, err := getBearerToken(
		client.Do(req),
	)
	if err != nil {
		return true, BearerTokenResponseError,
			fmt.Errorf("error in getting BearerToken response for image(%s): %v\n", imageName, err)
	}

	clairAddr, err := getAddress(kubeClient)
	if err != nil {
		return true, GettingClairAddrError,
			fmt.Errorf("error in getting ClairAddr for image(%s): %v", imageName, err)
	}
	for i := 0; i < layersLen; i++ {
		l := &LayerType{
			Name:       digest[7:] + layers[i].Digest[7:],
			Path:       fmt.Sprintf("%s/%s/%s/%s", registryUrl+"/v2", repo, "blobs", layers[i].Digest),
			ParentName: parent,
			Format:     "Docker",
			Headers: HeadersType{
				Authorization: token,
			},
		}
		parent = l.Name

		req, err := requestSendingLayer(l, clairAddr)
		if err != nil {
			return true, SendingLayerRequestError,
				fmt.Errorf("error in creating SendingLayerRequest for image(%s): %v\n", imageName, err)
		}
		_, err = clairClient.Do(req)
		if err != nil {
			return true, SendingLayerError,
				fmt.Errorf("error in sending layer of image(%s): %v\n", imageName, err)
		}
	}

	req, err = requestVulnerabilities(digest[7:]+layers[layersLen-1].Digest[7:], clairAddr)
	if err != nil {
		return true, VulnerabilitiesRequestError,
			fmt.Errorf("error in creating VulnerabilitiesRequest for image(%s): %v\n", imageName, err)
	}

	layerObj, err := decode(clairClient.Do(req))
	if err != nil {
		return true, VulnerabilitiesResponseError,
			fmt.Errorf("error in decoding VulnerabilitiesResponse for image(%s): %v\n", imageName, err)
	}

	//fs := getFeatures(layerObj)
	vuls := getVulnerabilities(layerObj)
	//oneliners.PrettyJson(fs, "Features")
	oneliners.PrettyJson(vuls, "vulnerabilities")
	if vuls != nil {
		return true, VulnerableStatus,
			fmt.Errorf("Image(%s) contains vulnerabilities", imageName)
	}

	return false, NotVulnerableStatus, nil
}
