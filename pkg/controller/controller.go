package controller

import (
	"net/http"

	"github.com/appscode/go/log"
	"github.com/appscode/pat"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type ScannerController struct {
	Config

	KubeClient kubernetes.Interface
	recorder   record.EventRecorder
}

func (c *ScannerController) RunOpsServer(stopCh <-chan struct{}) error {
	m := pat.New()
	m.Get("/metrics", promhttp.Handler())
	http.Handle("/", m)
	log.Infoln("Listening on", c.OpsAddress)
	return http.ListenAndServe(c.OpsAddress, nil)
}
