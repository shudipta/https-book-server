package controller

import (
	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/soter/scanner/pkg/clair"
	"github.com/soter/scanner/pkg/eventer"
	"github.com/soter/scanner/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type config struct {
	//ClairAddress  string
	//ClairCertDir  string
	FailurePolicy types.FailurePolicy
}

type Config struct {
	config

	ClientConfig   *rest.Config
	KubeClient     kubernetes.Interface
	WorkloadClient wcs.Interface
	Scanner        *clair.Scanner
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*Controller, error) {
	ctrl := &Controller{
		config: c.config,

		KubeClient:     c.KubeClient,
		WorkloadClient: c.WorkloadClient,
		recorder:       eventer.NewEventRecorder(c.KubeClient, "soter-scanner"),
		scanner:        c.Scanner,
	}
	return ctrl, nil
}
