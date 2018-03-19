package server

import (
	"flag"

	"github.com/soter/scanner/pkg/controller"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ControllerOptions struct {
	EnableRBAC bool
	OpsAddress string
}

func NewControllerOptions() *ControllerOptions {
	return &ControllerOptions{
		OpsAddress: ":56790",
	}
}

func (s *ControllerOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.OpsAddress, "ops-address", s.OpsAddress, "Address to listen on for web interface and telemetry.")
	fs.BoolVar(&s.EnableRBAC, "rbac", s.EnableRBAC, "Enable RBAC for operator")
}

func (s *ControllerOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("scanner", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ControllerOptions) ApplyTo(cfg *controller.ControllerConfig) error {
	var err error

	cfg.EnableRBAC = s.EnableRBAC
	cfg.OpsAddress = s.OpsAddress

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
