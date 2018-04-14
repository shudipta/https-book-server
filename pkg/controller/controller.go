package controller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type ScannerController struct {
	Config

	KubeClient kubernetes.Interface
	recorder   record.EventRecorder
}
