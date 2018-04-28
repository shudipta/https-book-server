package controller

import (
	utilerrors "github.com/appscode/go/util/errors"
	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) ScanCluster() error {
	var errs []error

	errs = append(errs, c.scanDeployments()...)
	errs = append(errs, c.scanReplicationControllers()...)
	errs = append(errs, c.scanReplicaSets()...)
	errs = append(errs, c.scanDaemonSets()...)
	errs = append(errs, c.scanJobs()...)
	errs = append(errs, c.scanCronJobs()...)
	errs = append(errs, c.scanStatefulSets()...)

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) scanDeployments() (errs []error) {
	objects, err := c.Client.AppsV1().Deployments(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanReplicationControllers() (errs []error) {
	objects, err := c.Client.CoreV1().ReplicationControllers(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanReplicaSets() (errs []error) {
	objects, err := c.Client.AppsV1().ReplicaSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanDaemonSets() (errs []error) {
	objects, err := c.Client.AppsV1().DaemonSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanJobs() (errs []error) {
	objects, err := c.Client.BatchV1().Jobs(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanCronJobs() (errs []error) {
	objects, err := c.Client.BatchV1beta1().CronJobs(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (c *Controller) scanStatefulSets() (errs []error) {
	objects, err := c.Client.AppsV1().StatefulSets(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errs = append(errs, err)
		return
	}

	for i := range objects.Items {
		if w, err := wcs.ConvertToWorkload(&objects.Items[i]); err != nil {
			errs = append(errs, err)
		} else {
			_, _, err := c.CheckWorkload(w)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}
