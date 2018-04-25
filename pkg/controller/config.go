package controller

import (
	"github.com/soter/scanner/pkg/eventer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type config struct {
	ClairAddress string
}

type Config struct {
	config

	ClientConfig *rest.Config
	KubeClient   kubernetes.Interface
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*Controller, error) {
	ctrl := &Controller{
		config: c.config,

		client:   c.KubeClient,
		recorder: eventer.NewEventRecorder(c.KubeClient, "scanner.soter.ac"),
	}
	return ctrl, nil
}
