package controller

import (
	hooks "github.com/appscode/kutil/admission/v1beta1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ScannerController) NewReplicationControllerWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "scanner.soter.cloud",
			Version:  "v1alpha1",
			Resource: "replicationcontrollers",
		},
		"replicationcontroller",
		core.SchemeGroupVersion.WithKind("ReplicationController"),
		nil,
		c,
	)
}
