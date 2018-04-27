package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	workload "github.com/appscode/kubernetes-webhook-util/apis/workload/v1"
	"github.com/appscode/kutil/tools/docker"
	_ "github.com/appscode/kutil/tools/docker"
	"github.com/coreos/clair/api/v3/clairpb"
	manifestV1 "github.com/docker/distribution/manifest/schema1"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"google.golang.org/grpc"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type Controller struct {
	config

	client   kubernetes.Interface
	recorder record.EventRecorder

	ClairAncestryServiceClient     clairpb.AncestryServiceClient
	ClairNotificationServiceClient clairpb.NotificationServiceClient
}

// ref: https://github.com/docker/cli/blob/master/vendor/github.com/docker/docker/api/types/auth.go
// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	// Email is an optional value associated with the username.
	// This field is deprecated and will be removed in a later
	// version of docker.
	Email string `json:"email,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identitytoken,omitempty"`

	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registrytoken,omitempty"`
}

func GetAllSecrets(refs []core.LocalObjectReference) []string {
	var names []string
	for _, ref := range refs {
		names = append(names, ref.Name)
	}
	return names
}

// Image represents Docker image
type Image struct {
	Ref  docker.ImageRef
	Auth *dockertypes.AuthConfig

	FsLayers      []string
	digest        string
	schemaVersion int
}

// checkContainers() checks vulnerabilities for each images used in containers.
// Here, precache parameter indicates that checking is being done for storing
// vulnerabilities and features of each image layer into cache. Otherwise,
// if precache is false then
// 		if any image is vulnerable then
//           this method returns
// This method takes namespace_name <namespace> of provided secrets <secretNames> and image name
// of a docker image. For each secret, it reads the config data of secret and store it to
// auth variable (map[string]map[string]AuthConfig)
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
func (c *Controller) CheckWorkload(w *workload.Workload) (*workload.Workload, bool, error) {
	var pullSecrets []core.Secret
	for _, ref := range w.Spec.Template.Spec.ImagePullSecrets {
		secret, err := c.client.CoreV1().Secrets(w.Namespace).Get(ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, false, err
		}
		pullSecrets = append(pullSecrets, *secret)
	}

	keyring, err := docker.MakeDockerKeyring(pullSecrets)
	if err != nil {
		return nil, false, err
	}

	images := sets.NewString()
	for _, c := range w.Spec.Template.Spec.Containers {
		images.Insert(c.Image)
	}
	for _, c := range w.Spec.Template.Spec.InitContainers {
		images.Insert(c.Image)
	}

	for _, image := range images.List() {
		ref, err := docker.ParseImageName(image)
		if err != nil {
			return nil, false, err
		}

		_, auth, mf, err := docker.PullManifest(ref, keyring)
		if err != nil {
			return nil, false, err
		}

		headers := map[string]string{}
		if auth.Username != "" && auth.Password != "" {
			headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth.Username+":"+auth.Password))
		}

		req := &clairpb.PostAncestryRequest{
			AncestryName: ref.String(),
			Format:       "Docker",
		}
		switch manifest := mf.(type) {
		case *manifestV2.DeserializedManifest:
			layers := make([]*clairpb.PostAncestryRequest_PostLayer, len(manifest.Layers))
			for i, layer := range manifest.Layers {
				layers[i] = &clairpb.PostAncestryRequest_PostLayer{
					Hash:    hashPart(manifest.Config.Digest.String()) + hashPart(layer.Digest.String()),
					Path:    fmt.Sprintf("%s/v2/%s/blobs/%s", auth.ServerAddress, ref.Repository, layer.Digest.String()),
					Headers: headers,
				}
			}
			req.Layers = layers
		case *manifestV1.SignedManifest:
			layers := make([]*clairpb.PostAncestryRequest_PostLayer, len(manifest.FSLayers))
			for i, layer := range manifest.FSLayers {
				layers[len(manifest.FSLayers)-1-i] = &clairpb.PostAncestryRequest_PostLayer{
					Hash:    hashPart(layer.BlobSum.String()),
					Path:    fmt.Sprintf("%s/v2/%s/blobs/%s", auth.ServerAddress, ref.Repository, layer.BlobSum.String()),
					Headers: headers,
				}
			}
			req.Layers = layers
		default:
			return nil, false, errors.New("unknown manifest type")
		}

		layersLen := len(req.Layers)
		if layersLen == 0 {
			return nil, false, errors.Wrapf(err, "failed to pull Layers for image %s", ref)
		}
		glog.Infoln("Analyzing", layersLen, "layers")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		_, err = c.ClairAncestryServiceClient.PostAncestry(ctx, req)
		if err != nil {
			return nil, false, errors.Wrapf(err, "failed to send layers for image %s", ref)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		resp, err := c.ClairAncestryServiceClient.GetAncestry(ctx, &clairpb.GetAncestryRequest{
			//AncestryName:        ref.Repository,
			AncestryName:        ref.String(),
			WithFeatures:        true,
			WithVulnerabilities: true,
		})
		if err != nil {
			return nil, false, err
		}

		vulnerabilities := GetVulnerabilities(resp)
		if len(vulnerabilities) > 0 {
			return nil, false, fmt.Errorf("Image %s contains vulnerabilities", ref)
		}
		// if resp.Ancestry.Features
		// 	return getFeatures(resp), getVulnerabilities(resp), nil
	}
	return w, true, nil
}

func hashPart(digest string) string {
	if len(digest) < 7 {
		return ""
	}

	return digest[7:]
}

func GetVulnerabilities(resp *clairpb.GetAncestryResponse) []api.Vulnerability {
	var vuls []api.Vulnerability
	if resp == nil {
		return nil
	}
	for _, feature := range resp.Ancestry.Features {
		for _, vul := range feature.Vulnerabilities {
			vuls = append(vuls, api.Vulnerability{
				Name:          vul.Name,
				NamespaceName: vul.NamespaceName,
				Description:   vul.Description,
				Link:          vul.Link,
				Severity:      vul.Severity,
				FixedBy:       vul.FixedBy,
				FeatureName:   feature.Name,
			})
		}
	}

	return vuls
}

func getFeatures(resp *clairpb.GetAncestryResponse) []api.Feature {
	var fs []api.Feature
	if resp == nil {
		return nil
	}
	for _, feature := range resp.Ancestry.Features {
		fs = append(fs, api.Feature{
			Name:          feature.Name,
			NamespaceName: feature.NamespaceName,
			Version:       feature.Version,
		})
	}

	return fs
}

func clairClientSetup(clairAddress string) (clairpb.AncestryServiceClient, error) {
	conn, err := grpc.Dial(clairAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := clairpb.NewAncestryServiceClient(conn)
	return c, nil
}
