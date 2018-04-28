package controller

import (
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	"github.com/appscode/kubernetes-webhook-util/admission/v1beta1/workload"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Controller) NewDeploymentWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "deployments",
		},
		"deployment",
		"Deployment",
		nil,
		c,
	)
}

func (c *Controller) NewReplicaSetWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "replicasets",
		},
		"replicaset",
		"ReplicaSet",
		nil,
		c,
	)
}

func (c *Controller) NewDaemonSetWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "daemonsets",
		},
		"daemonset",
		"DaemonSet",
		nil,
		c,
	)
}

func (c *Controller) NewStatefulSetWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "statefulsets",
		},
		"statefulset",
		"StatefulSet",
		nil,
		c,
	)
}

func (c *Controller) NewReplicationControllerWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "replicationcontrollers",
		},
		"replicationcontroller",
		"ReplicationController",
		nil,
		c,
	)
}

func (c *Controller) NewJobWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "jobs",
		},
		"job",
		"Job",
		nil,
		c,
	)
}

func (c *Controller) NewCronJobWebhook() hooks.AdmissionHook {
	return workload.NewWorkloadWebhook(
		schema.GroupVersionResource{
			Group:    "admission.scanner.soter.ac",
			Version:  "v1alpha1",
			Resource: "cronjobs",
		},
		"cronjob",
		"CronJob",
		nil,
		c,
	)
}
