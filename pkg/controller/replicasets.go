package controller

import (
	"github.com/appscode/kutil/admission"
	hooks "github.com/appscode/kutil/admission/v1beta1"
	workload "github.com/appscode/kutil/workload/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func (c *ScannerController) NewReplicaSetWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.soter.cloud",
			Version:  "v1alpha1",
			Resource: "replicasets",
		},
		"replicaset",
		apps.SchemeGroupVersion.WithKind("ReplicaSet"),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkReplicaSet(obj.(*workload.Workload))
				return modObj, err

			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				modObj, _, err := c.checkReplicaSet(newObj.(*workload.Workload))
				return modObj, err
			},
		},
	)
}

func (c *ScannerController) checkReplicaSet(w *workload.Workload) (*workload.Workload, bool, error) {
	return w, false, nil
}
