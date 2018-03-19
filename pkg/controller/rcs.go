package controller

import (
	"github.com/appscode/kutil/admission"
	hooks "github.com/appscode/kutil/admission/v1beta1"
	workload "github.com/appscode/kutil/workload/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ScannerController) NewReplicationControllerWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.soter.cloud",
			Version:  "v1alpha1",
			Resource: "replicationcontrollers",
		},
		"replicationcontroller",
		core.SchemeGroupVersion.WithKind("ReplicationController"),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkReplicationController(obj.(*workload.Workload))
				return modObj, err

			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkReplicationController(newObj.(*workload.Workload))
				return modObj, err
			},
		},
	)
}

func (c *ScannerController) checkReplicationController(w *workload.Workload) (*workload.Workload, bool, error) {
	return w, false, nil
}
