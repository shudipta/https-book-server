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
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
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
	ParseImageNameError          = 3
	GettingManifestError         = 4
	GettingCannonicalError       = 5
	DecodingCannonical_1_Error   = 6
	DecodingCannonical_2_Error   = 7
	PullingLayersError           = 8
	BearerTokenRequestError      = 9
	BearerTokenResponseError     = 10
	GettingClairAddrError        = 11
	SendingLayerRequestError     = 12
	SendingLayerError            = 13
	VulnerabilitiesRequestError  = 14
	VulnerabilitiesResponseError = 15
	VulnerableStatus             = 16
	NotVulnerableStatus          = 17
)

// This method takes <registryUrl>, <imageName>, <username>, <password> and
// returns (true, error) if any error occurred or vulnerability found otherwise return
// (false, <nil>). First, it parses the <imageName> and connect to the registry. Then,
// it gets the manifests for the <imageName>. Then for each layer in the manifests it
// sends a LayerType{} obj to clair. Then it does a GET requests to clair and receives
// a LayerType{} obj as response.Body. Then it filters this layer obj to get the features
// and vulnerabilities. Finally, it stores them in cache against the layer name as key.
// It stores them so that it finds them in cache without calling to clair and filtering
// again if next time we need to scan same layer.
// For more information about LayerType{}, https://coreos.com/clair/docs/latest/api_v1.html
// will be helpful.
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

	// TODO: need to check for digest part
	registryUrl, repo, tag, _, err := parseImageName(imageName, registryUrl)
	if err != nil {
		return Canonical1{}, ParseImageNameError, err
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
			fs   = []api.Feature{}
			vuls = []api.Vulnerability{}
		)

		if val, exists := vulsCache.Get(lName); exists {
			vuls = val.([]api.Vulnerability)
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
			}
			fs = getFeatures(layerObj)
			vuls = getVulnerabilities(layerObj)

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
