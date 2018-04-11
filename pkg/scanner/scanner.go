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
	GettingManifestError         = 3
	GettingCannonicalError       = 4
	DecodingCannonical_1_Error   = 5
	DecodingCannonical_2_Error   = 6
	PullingLayersError           = 7
	BearerTokenRequestError      = 8
	BearerTokenResponseError     = 9
	GettingClairAddrError        = 10
	SendingLayerRequestError     = 11
	SendingLayerError            = 12
	VulnerabilitiesRequestError  = 13
	VulnerabilitiesResponseError = 14
	VulnerableStatus             = 15
	NotVulnerableStatus          = 16
)

func IsVulnerable(
	kc kubernetes.Interface, fsCache, vulsCache *lru.TwoQueueCache,
	registryUrl, imageName, username, password string,
	precache bool) (Canonical1, int, error) {

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
		return Canonical1{}, GettingManifestError,
			fmt.Errorf("error in getting manifest for image(%s): %v\n", imageName, err)
	}
	canonicalBytes, err := manifest.MarshalJSON()
	if err != nil {
		return Canonical1{}, GettingCannonicalError,
			fmt.Errorf("error in getting manifest.canonical for image(%s): %v\n", imageName, err)
	}

	var imageManifest Canonical1
	if err := json.NewDecoder(bytes.NewReader(canonicalBytes)).Decode(&imageManifest); err != nil {
		return Canonical1{}, DecodingCannonical_1_Error,
			fmt.Errorf("error in decoding into canonical1 for image(%s): %v\n", imageName, err)
	}
	oneliners.PrettyJson(imageManifest, "image")

	if imageManifest.Layers == nil {
		var image2 Canonical2
		if err := json.NewDecoder(bytes.NewReader(canonicalBytes)).Decode(&image2); err != nil {
			return Canonical1{}, DecodingCannonical_2_Error,
				fmt.Errorf("error in decoding into canonical2 for image(%s): %v\n", imageName, err)
		}

		imageManifest.Layers = make([]layer1, len(image2.FsLayers))
		for i, l := range image2.FsLayers {
			imageManifest.Layers[len(image2.FsLayers)-1-i].Digest = l.BlobSum
		}
		imageManifest.SchemaVersion = image2.SchemaVersion
		//imageManifest.Config.Digest = "sha256:1df6g5874tryerjn6549a8d461vs6rf41468astgretv41xcb54ser6"
	}

	layersLen := len(imageManifest.Layers)
	if layersLen == 0 {
		return imageManifest, PullingLayersError,
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
		return imageManifest, BearerTokenRequestError,
			fmt.Errorf("error in creating BearerToken request for image(%s): %v\n", imageName, err)
	}

	token, err := getBearerToken(
		client.Do(req),
	)
	if err != nil {
		return imageManifest, BearerTokenResponseError,
			fmt.Errorf("error in getting BearerToken response for image(%s): %v\n", imageName, err)
	}

	clairAddr, err := getAddress(kc)
	if err != nil {
		return imageManifest, GettingClairAddrError,
			fmt.Errorf("error in getting ClairAddr for image(%s): %v", imageName, err)
	}

	for i := 0; i < layersLen; i++ {
		lName := HashPart(imageManifest.Config.Digest) + HashPart(imageManifest.Layers[i].Digest)

		var (
			fs   = []Feature{}
			vuls = []Vulnerability{}
		)

		if val, exists := vulsCache.Get(lName); exists {
			vuls = val.([]Vulnerability)
		} else {
			l := &LayerType{
				Name: lName,
				Path: fmt.Sprintf("%s/%s/%s/%s", registryUrl+"/v2", repo, "blobs", imageManifest.Layers[i].Digest),
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
				return imageManifest, SendingLayerRequestError,
					fmt.Errorf("error in creating SendingLayerRequest for image(%s).layer[%d]: %v\n", imageName, i, err)
			}
			_, err = clairClient.Do(req)
			if err != nil {
				return imageManifest, SendingLayerError,
					fmt.Errorf("error in sending layer of image(%s).layer[%d]: %v\n", imageName, i, err)
			}

			req, err = requestVulnerabilities(l.Name, clairAddr)
			if err != nil {
				return imageManifest, VulnerabilitiesRequestError,
					fmt.Errorf("error in creating VulnerabilitiesRequest for image(%s).layer[%d]: %v\n", imageName, i, err)
			}

			layerObj, err := decode(clairClient.Do(req))
			if err != nil {
				return imageManifest, VulnerabilitiesResponseError,
					fmt.Errorf("error in decoding VulnerabilitiesResponse for image(%s).layer[%d]: %v\n", imageName, i, err)
			} else {
				fs = getFeatures(layerObj)
				vuls = getVulnerabilities(layerObj)
			}

			//oneliners.PrettyJson(fs, "Features")
			//oneliners.PrettyJson(vuls, "vulnerabilities")
			cacheFeaturesAndVulnerabilities(fsCache, vulsCache, l.Name, fs, vuls)
		}
		if !precache && vuls != nil {
			return imageManifest, VulnerableStatus,
				fmt.Errorf("Image(%s) contains vulnerabilities", imageName)
		}
	}

	return imageManifest, NotVulnerableStatus, nil
}
