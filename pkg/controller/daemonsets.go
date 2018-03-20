package controller

import (
	hooks "github.com/appscode/kutil/admission/v1beta1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ScannerController) NewDaemonSetWebhook() hooks.AdmissionHook {
	return hooks.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "scanner.soter.cloud",
			Version:  "v1alpha1",
			Resource: "daemonsets",
		},
		"daemonset",
		extensions.SchemeGroupVersion.WithKind("DaemonSet"),
		nil,
		c,
	)
}
