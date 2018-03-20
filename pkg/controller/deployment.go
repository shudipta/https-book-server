package controller

import (
	hooks "github.com/appscode/kutil/admission/v1beta1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ScannerController) NewDeploymentWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "scanner.soter.cloud",
			Version:  "v1alpha1",
			Resource: "deployments",
		},
		"deployment",
		appsv1beta1.SchemeGroupVersion.WithKind("Deployment"),
		nil,
		c,
	)
}
