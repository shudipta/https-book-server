package controller

import (
	hooks "github.com/appscode/kutil/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func (c *ScannerController) NewJobWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "scanner.soter.cloud",
			Version:  "v1alpha1",
			Resource: "jobs",
		},
		"job",
		apps.SchemeGroupVersion.WithKind("Job"),
		nil,
		c,
	)
}
