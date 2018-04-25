package ping

import (
	"sync"

	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"github.com/soter/scanner/client/clientset/versioned"
	"github.com/soter/scanner/pkg/clair"
	"github.com/soter/scanner/pkg/controller"
	"github.com/tamalsaha/go-oneliners"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	restconfig "k8s.io/client-go/rest"
)

type REST struct {
	client     versioned.Interface
	controller *controller.Controller

	lock sync.RWMutex
}

var _ rest.Creater = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(config *restconfig.Config, ctl *controller.Controller) *REST {
	return &REST{
		client:     versioned.NewForConfigOrDie(config),
		controller: ctl,
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
	ns := apirequest.NamespaceValue(ctx)
	secretNames := controller.GetAllSecrets(req.Request.ImagePullSecrets)

	features, vulnerabilities, err := r.controller.CheckImage(ns, req.Request.Image, secretNames)
	if err.(*clair.ErrorWithCode).Code() != clair.VulnerableStatus &&
		err.(*clair.ErrorWithCode).Code() != clair.NotVulnerableStatus {
		return nil, err
	}

	req.Response = &api.ImageReviewResponse{
		Features:        features,
		Vulnerabilities: vulnerabilities,
	}

	oneliners.PrettyJson(req.Response)

	return req, nil
}
