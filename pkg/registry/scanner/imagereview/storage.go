package ping

import (
	_ "github.com/hashicorp/golang-lru"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"github.com/soter/scanner/pkg/clair"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

type REST struct {
	scanner *clair.Scanner
	kind    string
}

var _ rest.Creater = &REST{}
var _ rest.Getter = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(scanner *clair.Scanner, kind string) *REST {
	return &REST{
		scanner: scanner,
		kind:    kind,
	}
}

func (r *REST) New() runtime.Object {
	return &api.ImageReview{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(r.kind)
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*api.ImageReview)
	namespace := apirequest.NamespaceValue(ctx)

	result, err := r.scanner.ScanWorkload(r.kind, req.Name, namespace)
	if err != nil {
		return nil, err
	}
	req.Response = &api.ImageReviewResponse{Images: result}
	return req, nil
}

func (r *REST) Get(ctx apirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := apirequest.NamespaceValue(ctx)

	result, err := r.scanner.ScanWorkload(r.kind, name, namespace)
	if err != nil {
		return nil, err
	}
	return &api.ImageReview{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: apirequest.NamespaceValue(ctx),
		},
		Response: &api.ImageReviewResponse{
			Images: result,
		},
	}, nil
}
