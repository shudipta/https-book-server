package root

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *InitSetter) scanDeployments(namespace string) {
	deploys, err := s.ctl.KubeClient.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, deploy := range deploys.Items {
		s.ctl.CheckWorkload(convertToWorkload(&deploy), true)
	}
}

func (s *InitSetter) scanReplicationControllers(namespace string) {
	rcs, err := s.ctl.KubeClient.CoreV1().ReplicationControllers(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, rc := range rcs.Items {
		s.ctl.CheckWorkload(convertToWorkload(&rc), true)
	}
}

func (s *InitSetter) scanReplicaSets(namespace string) {
	rss, err := s.ctl.KubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, rs := range rss.Items {
		s.ctl.CheckWorkload(convertToWorkload(&rs), true)
	}
}

func (s *InitSetter) scanDaemonSets(namespace string) {
	dss, err := s.ctl.KubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, ds := range dss.Items {
		s.ctl.CheckWorkload(convertToWorkload(&ds), true)
	}
}

func (s *InitSetter) scanJobs(namespace string) {
	jobs, err := s.ctl.KubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, job := range jobs.Items {
		s.ctl.CheckWorkload(convertToWorkload(&job), true)
	}
}

func (s *InitSetter) scanStatefulSets(namespace string) {
	stss, err := s.ctl.KubeClient.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		errorForGettingWorkloadFromNamespac(namespace, err)
	}

	for _, sts := range stss.Items {
		s.ctl.CheckWorkload(convertToWorkload(&sts), true)
	}
}
