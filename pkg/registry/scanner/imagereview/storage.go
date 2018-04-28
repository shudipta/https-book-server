package ping

import (
	"context"
	"sync"
	"time"

	"github.com/appscode/kutil/tools/docker"
	"github.com/coreos/clair/api/v3/clairpb"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"github.com/soter/scanner/client/clientset/versioned"
	"github.com/soter/scanner/pkg/controller"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	restconfig "k8s.io/client-go/rest"
)

type REST struct {
	client     versioned.Interface
	controller *controller.Controller
	imageRefs  map[string]string

	lock sync.RWMutex
}

var _ rest.Getter = &REST{}
var _ rest.Creater = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(config *restconfig.Config, ctl *controller.Controller) *REST {
	return &REST{
		client:     versioned.NewForConfigOrDie(config),
		controller: ctl,
		imageRefs:  make(map[string]string),
	}
}

func (r *REST) New() runtime.Object {
	return &api.ImageReview{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(api.ResourceKindImageReview)
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*api.ImageReview)
	namespace := apirequest.NamespaceValue(ctx)
	var pullSecrets []core.Secret
	for _, ref := range req.Request.ImagePullSecrets {
		secret, err := r.controller.KubeClient.CoreV1().Secrets(namespace).Get(ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		pullSecrets = append(pullSecrets, *secret)
	}

	keyring, err := docker.MakeDockerKeyring(pullSecrets)
	if err != nil {
		return nil, err
	}

	ref, err := docker.ParseImageName(req.Request.Image)
	if err != nil {
		return nil, err
	}
	r.imageRefs[req.Name] = ref.String()

	_, auth, mf, err := docker.PullManifest(ref, keyring)
	if err != nil {
		return nil, err
	}

	postReq, err := r.controller.NewPostAncestryRequest(ref, auth, mf)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err = r.controller.ClairAncestryServiceClient.PostAncestry(ctx, postReq)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *REST) Get(ctx apirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	resp, err := r.controller.ClairAncestryServiceClient.GetAncestry(ctx, &clairpb.GetAncestryRequest{
		AncestryName:        r.imageRefs[name],
		WithFeatures:        true,
		WithVulnerabilities: true,
	})
	if err != nil {
		return nil, err
	}

	delete(r.imageRefs, name)

	features := getFeatures(resp)
	vulnerabilities := getVulnerabilities(resp)

	return &api.ImageReview{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: apirequest.NamespaceValue(ctx),
		},
		Response: &api.ImageReviewResponse{
			Features:        features,
			Vulnerabilities: vulnerabilities,
		},
	}, nil
}

func getVulnerabilities(resp *clairpb.GetAncestryResponse) []api.Vulnerability {
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
