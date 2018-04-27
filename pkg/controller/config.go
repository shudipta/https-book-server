package controller

import (
	"fmt"

	"github.com/soter/scanner/pkg/clair-api"
	"github.com/soter/scanner/pkg/eventer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type config struct {
	ClairAddress    string
	ClairApiCertDir string
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
	dialOption, err := clair_api.DialOptionForTLSConfig(c.ClairApiCertDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial option for tls: %v", err)
	}

	clairAncestryServiceClient, err := clair_api.NewClairAncestryServiceClient(c.ClairAddress, dialOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for Ancestry Service: %v", err)
	}

	clairNotificationServiceClient, err := clair_api.NewClairNotificationServiceClient(c.ClairAddress, dialOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for Notification Service: %v", err)
	}

	ctrl := &Controller{
		config: c.config,

		client:   c.KubeClient,
		recorder: eventer.NewEventRecorder(c.KubeClient, "soter-scanner"),

		ClairAncestryServiceClient:     clairAncestryServiceClient,
		ClairNotificationServiceClient: clairNotificationServiceClient,
	}
	return ctrl, nil
}
