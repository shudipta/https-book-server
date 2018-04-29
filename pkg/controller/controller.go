package controller

import (
	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/soter/scanner/pkg/clair"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type Controller struct {
	config

	KubeClient     kubernetes.Interface
	WorkloadClient wcs.Interface
	recorder       record.EventRecorder
	scanner        *clair.Scanner
}

func (c *Controller) ScanCluster() error {
	return c.scanner.ScanCluster()
}
