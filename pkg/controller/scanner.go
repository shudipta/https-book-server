package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	utilerrors "github.com/appscode/go/util/errors"
	wpi "github.com/appscode/kubernetes-webhook-util/apis/workload/v1"
	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/appscode/kutil/tools/docker"
	_ "github.com/appscode/kutil/tools/docker"
	"github.com/coreos/clair/api/v3/clairpb"
	manifestV1 "github.com/docker/distribution/manifest/schema1"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/soter/scanner/pkg/types"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type Controller struct {
	config

	KubeClient     kubernetes.Interface
	WorkloadClient wcs.Interface
	recorder       record.EventRecorder

	ClairAncestryServiceClient     clairpb.AncestryServiceClient
	ClairNotificationServiceClient clairpb.NotificationServiceClient
}

func (c *Controller) ScanCluster() error {
	var errs []error

	objects, err := c.WorkloadClient.Workloads(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
	} else {
		for i := range objects.Items {
			if err := c.checkWorkload(&objects.Items[i]); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
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
func (c *Controller) checkWorkload(w *wpi.Workload) error {
	var pullSecrets []core.Secret
	for _, ref := range w.Spec.Template.Spec.ImagePullSecrets {
		secret, err := c.KubeClient.CoreV1().Secrets(w.Namespace).Get(ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		pullSecrets = append(pullSecrets, *secret)
	}

	keyring, err := docker.MakeDockerKeyring(pullSecrets)
	if err != nil {
		return err
	}

	images := sets.NewString()
	for _, c := range w.Spec.Template.Spec.Containers {
		images.Insert(c.Image)
	}
	for _, c := range w.Spec.Template.Spec.InitContainers {
		images.Insert(c.Image)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, image := range images.List() {
		ref, err := docker.ParseImageName(image)
		if err != nil {
			return err
		}
		_, auth, mf, err := docker.PullManifest(ref, keyring)
		if err != nil {
			if c.FailurePolicy == types.FailurePolicyIgnore {
				continue
			}
			return err
		}

		req, err := c.NewPostAncestryRequest(ref, auth, mf)
		if err != nil {
			return err
		}

		_, err = c.ClairAncestryServiceClient.PostAncestry(ctx, req)
		if err != nil {
			return errors.Wrapf(err, "failed to send layers for image %s", ref)
		}

		resp, err := c.ClairAncestryServiceClient.GetAncestry(ctx, &clairpb.GetAncestryRequest{
			AncestryName:        ref.String(),
			WithFeatures:        true,
			WithVulnerabilities: true,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to get layers for image %s", ref)
		}

		if hasVulnerabilities(resp) {
			return errors.Errorf("image %s contains vulnerabilities", ref)
		}
	}

	return nil
}

func (c *Controller) NewPostAncestryRequest(ref docker.ImageRef, auth *dockertypes.AuthConfig, mf interface{}) (*clairpb.PostAncestryRequest, error) {
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
		return nil, errors.New("unknown manifest type")
	}
	if len(req.Layers) == 0 {
		return nil, errors.Errorf("failed to pull Layers for image %s", ref)
	}
	return req, nil
}

func hashPart(digest string) string {
	if len(digest) < 7 {
		return ""
	}

	return digest[7:]
}

func hasVulnerabilities(resp *clairpb.GetAncestryResponse) bool {
	if resp == nil {
		return false
	}
	for _, feature := range resp.Ancestry.Features {
		if len(feature.Vulnerabilities) > 0 {
			return true
		}
	}

	return false
}
