package server

import (
	"flag"

	"github.com/soter/scanner/pkg/controller"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ControllerOptions struct {
	EnableRBAC bool
}

func NewControllerOptions() *ControllerOptions {
	return &ControllerOptions{}
}

func (s *ControllerOptions) AddGoFlags(fs *flag.FlagSet) {
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

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
