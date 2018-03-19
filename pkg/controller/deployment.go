package controller

import (
	"github.com/appscode/kutil/admission"
	hooks "github.com/appscode/kutil/admission/v1beta1"
	workload "github.com/appscode/kutil/workload/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ScannerController) NewDeploymentWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.soter.cloud",
			Version:  "v1alpha1",
			Resource: "deployments",
		},
		"deployment",
		appsv1beta1.SchemeGroupVersion.WithKind("Deployment"),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkDeployment(obj.(*workload.Workload))
				return modObj, err

			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkDeployment(newObj.(*workload.Workload))
				return modObj, err
			},
		},
	)
}

func (c *ScannerController) checkDeployment(w *workload.Workload) (*workload.Workload, bool, error) {
	return w, false, nil
}
