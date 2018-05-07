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
	"github.com/tamalsaha/go-oneliners"
)

type ScannerStorage interface {
	rest.Storage
	Resource() (plural schema.GroupVersionResource, singular string)
}

type REST struct {
	scanner  *clair.Scanner
	plural   schema.GroupVersionResource
	singular string
}

var _ rest.Creater = &REST{}
var _ rest.Getter = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(scanner *clair.Scanner, plural schema.GroupVersionResource, singular string) *REST {
	return &REST{
		scanner:  scanner,
		plural:   plural,
		singular: singular,
	}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind("ImageReview")
}

func (r *REST) Resource() (schema.GroupVersionResource, string) {
	return r.plural, r.singular
}

func (r *REST) New() runtime.Object {
	return &api.ImageReview{}
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*api.ImageReview)
	oneliners.PrettyJson(r, "rest for "+req.Name)
	oneliners.PrettyJson(req, "rest create req")
	namespace := apirequest.NamespaceValue(ctx)

	result, err := r.scanner.ScanWorkload(r.singular, req.Name, namespace)
	if err != nil {
		return nil, err
	}
	req.Response = &api.ImageReviewResponse{Images: result}
	return req, nil
}

func (r *REST) Get(ctx apirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := apirequest.NamespaceValue(ctx)

	result, err := r.scanner.ScanWorkload(r.singular, name, namespace)
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
