package scanner

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru"
	reg "github.com/heroku/docker-registry-client/registry"
	"github.com/tamalsaha/go-oneliners"
	"k8s.io/client-go/kubernetes"
)

type config struct {
	MediaType string
	Size      int
	Digest    string
}

type layer1 struct {
	MediaType string
	Size      int
	Digest    string
}

type Canonical1 struct {
	SchemaVersion int
	MediaType     string
	Config        config
	Layers        []layer1
}

type Canonical2 struct {
	SchemaVersion int
	FsLayers      []layer2
}

type layer2 struct {
	BlobSum string
}

const (
	GettingManifestError         = 1
	GettingCannonicalError       = 2
	DecodingCannonical_1_Error   = 3
	DecodingCannonical_2_Error   = 4
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

func IsVulnerable(
	kc kubernetes.Interface, fsCache, vulsCache *lru.TwoQueueCache,
	registryUrl, imageName, username, password string,
	precache bool) (bool, int, error) {

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, //true
			},
		},
		Timeout: time.Minute,
	}

	registry, repo, tag := parseImageName(imageName)
	if registryUrl == "" {
		registryUrl = registry
	}

	hub := &reg.Registry{
		URL: registryUrl,
		Client: &http.Client{
			Transport: reg.WrapTransport(http.DefaultTransport, registryUrl, username, password),
		},
		Logf: reg.Quiet,
	}

	manifest, err := hub.ManifestV2(repo, tag)
	if err != nil {
		return true, GettingManifestError,
			fmt.Errorf("error in getting manifest for image(%s): %v\n", imageName, err)
	}
	canonicalBytes, err := manifest.MarshalJSON()
	if err != nil {
		return true, GettingCannonicalError,
			fmt.Errorf("error in getting manifest.canonical for image(%s): %v\n", imageName, err)
	}

	var image Canonical1
	if err := json.NewDecoder(bytes.NewReader(canonicalBytes)).Decode(&image); err != nil {
		return true, DecodingCannonical_1_Error,
			fmt.Errorf("error in decoding into canonical1 for image(%s): %v\n", imageName, err)
	}
	oneliners.PrettyJson(image, "image")

	if image.Layers == nil {
		var image2 Canonical2
		if err := json.NewDecoder(bytes.NewReader(canonicalBytes)).Decode(&image2); err != nil {
			return true, DecodingCannonical_2_Error,
				fmt.Errorf("error in decoding into canonical2 for image(%s): %v\n", imageName, err)
		}

		image.Layers = make([]layer1, len(image2.FsLayers))
		for i, l := range image2.FsLayers {
			image.Layers[len(image2.FsLayers)-1-i].Digest = l.BlobSum
		}
		image.SchemaVersion = image2.SchemaVersion
	}

	digest := image.Config.Digest

	layersLen := len(image.Layers)
	if layersLen == 0 {
		return true, PullingLayersError,
			fmt.Errorf("error is pulling fsLayers for image(%s)\n", imageName)
	} else {
		fmt.Println("Analysing", layersLen, "layers")
	}

	clairClient := http.Client{
		Timeout: time.Minute,
	}

	//var parent string
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

	clairAddr, err := getAddress(kc)
	if err != nil {
		return true, GettingClairAddrError,
			fmt.Errorf("error in getting ClairAddr for image(%s): %v", imageName, err)
	}

	for i := 0; i < layersLen; i++ {
		lName := hashPart(digest) + hashPart(image.Layers[i].Digest)

		var (
			fs   []Feature       = nil
			vuls []Vulnerability = nil
		)

		if val, exists := vulsCache.Get(lName); exists {
			vuls = val.([]Vulnerability)
		} else {
			l := &LayerType{
				Name: lName,
				Path: fmt.Sprintf("%s/%s/%s/%s", registryUrl+"/v2", repo, "blobs", image.Layers[i].Digest),
				//ParentName: parent,
				ParentName: "",
				Format:     "Docker",
				Headers: HeadersType{
					Authorization: token,
				},
			}
			//parent = l.Name

			req, err := requestSendingLayer(l, clairAddr)
			if err != nil {
				return true, SendingLayerRequestError,
					fmt.Errorf("error in creating SendingLayerRequest for image(%s).layer[%d]: %v\n", imageName, i, err)
			}
			_, err = clairClient.Do(req)
			if err != nil {
				return true, SendingLayerError,
					fmt.Errorf("error in sending layer of image(%s).layer[%d]: %v\n", imageName, i, err)
			}

			req, err = requestVulnerabilities(l.Name, clairAddr)
			if err != nil {
				return true, VulnerabilitiesRequestError,
					fmt.Errorf("error in creating VulnerabilitiesRequest for image(%s).layer[%d]: %v\n", imageName, i, err)
			}

			layerObj, err := decode(clairClient.Do(req))
			if err != nil {
				return true, VulnerabilitiesResponseError,
					fmt.Errorf("error in decoding VulnerabilitiesResponse for image(%s).layer[%d]: %v\n", imageName, i, err)
			}

			fs = getFeatures(layerObj)
			vuls = getVulnerabilities(layerObj)
			//oneliners.PrettyJson(fs, "Features")
			//oneliners.PrettyJson(vuls, "vulnerabilities")
			cacheFeaturesAndVulnerabilities(fsCache, vulsCache, l.Name, fs, vuls)
		}
		if !precache && vuls != nil {
			return true, VulnerableStatus,
				fmt.Errorf("Image(%s) contains vulnerabilities", imageName)
		}
	}

	return false, NotVulnerableStatus, nil
}
