package controller

import (
	"net/http"

	"github.com/appscode/go/log"
	"github.com/appscode/pat"
	"github.com/hashicorp/golang-lru"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type ScannerController struct {
	Config

	KubeClient kubernetes.Interface
	recorder   record.EventRecorder

	FsCache   *lru.TwoQueueCache
	VulsCache *lru.TwoQueueCache
}

func (c *ScannerController) RunOpsServer(stopCh <-chan struct{}) error {
	//cache.New(c).Run()

	m := pat.New()
	m.Get("/metrics", promhttp.Handler())
	http.Handle("/", m)
	log.Infoln("Listening on", c.OpsAddress)
	return http.ListenAndServe(c.OpsAddress, nil)
}
