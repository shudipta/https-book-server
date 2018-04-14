package cache

import (
	workload "github.com/appscode/kubernetes-webhook-util/workload/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *PreScanner) scanDeployments() (errs []error) {
	objects, err := s.scanner.KubeClient.AppsV1().Deployments(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanReplicationControllers() (errs []error) {
	objects, err := s.scanner.KubeClient.CoreV1().ReplicationControllers(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanReplicaSets() (errs []error) {
	objects, err := s.scanner.KubeClient.AppsV1().ReplicaSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanDaemonSets() (errs []error) {
	objects, err := s.scanner.KubeClient.AppsV1().DaemonSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanJobs() (errs []error) {
	objects, err := s.scanner.KubeClient.BatchV1().Jobs(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanCronJobs() (errs []error) {
	objects, err := s.scanner.KubeClient.BatchV1beta1().CronJobs(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *PreScanner) scanStatefulSets() (errs []error) {
	objects, err := s.scanner.KubeClient.AppsV1().StatefulSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := workload.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := s.scanner.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}
