package controller

import (
	"time"

	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	"github.com/hashicorp/golang-lru"
	"github.com/soter/scanner/pkg/eventer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	EnableRBAC      bool
	ScannerImageTag string
	DockerRegistry  string
	MaxNumRequeues  int
	NumThreads      int
	OpsAddress      string
	ResyncPeriod    time.Duration
}

type ControllerConfig struct {
	Config

	ClientConfig   *rest.Config
	KubeClient     kubernetes.Interface
	AdmissionHooks []hooks.AdmissionHook
}

func NewControllerConfig(clientConfig *rest.Config) *ControllerConfig {
	return &ControllerConfig{
		ClientConfig: clientConfig,
	}
}

func (c *ControllerConfig) New() (*ScannerController, error) {
	fsCache, err := lru.New2Q(128)
	if err != nil {
		return nil, err
	}
	vulsCache, err := lru.New2Q(128)
	if err != nil {
		return nil, err
	}

	ctrl := &ScannerController{
		Config:     c.Config,
		KubeClient: c.KubeClient,
		FsCache:    fsCache,
		VulsCache:  vulsCache,
		recorder:   eventer.NewEventRecorder(c.KubeClient, "soter-scanner"),
	}
	return ctrl, nil
}
