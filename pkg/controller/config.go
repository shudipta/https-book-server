package controller

import (
	"fmt"

	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/soter/scanner/pkg/clair"
	"github.com/soter/scanner/pkg/eventer"
	"github.com/soter/scanner/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type config struct {
	ClairAddress    string
	ClairApiCertDir string
	FailurePolicy   types.FailurePolicy
}

type Config struct {
	config

	ClientConfig   *rest.Config
	KubeClient     kubernetes.Interface
	WorkloadClient wcs.Interface
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*Controller, error) {
	dialOption, err := clair.DialOptionForTLSConfig(c.ClairApiCertDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial option for tls: %v", err)
	}

	clairAncestryServiceClient, err := clair.NewClairAncestryServiceClient(c.ClairAddress, dialOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for Ancestry Service: %v", err)
	}

	clairNotificationServiceClient, err := clair.NewClairNotificationServiceClient(c.ClairAddress, dialOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for Notification Service: %v", err)
	}

	ctrl := &Controller{
		config: c.config,

		KubeClient:     c.KubeClient,
		WorkloadClient: c.WorkloadClient,
		recorder:       eventer.NewEventRecorder(c.KubeClient, "soter-scanner"),

		ClairAncestryServiceClient:     clairAncestryServiceClient,
		ClairNotificationServiceClient: clairNotificationServiceClient,
	}

	return ctrl, nil
}
