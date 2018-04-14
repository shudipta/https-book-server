package server

import (
	"flag"

	"github.com/soter/scanner/pkg/controller"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ControllerOptions struct {
	QPS        float64
	Burst      int
	EnableRBAC bool
}

func NewControllerOptions() *ControllerOptions {
	return &ControllerOptions{
		QPS:   100,
		Burst: 100,
	}
}

func (s *ControllerOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")

	fs.BoolVar(&s.EnableRBAC, "rbac", s.EnableRBAC, "Enable RBAC for operator")
}

func (s *ControllerOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("scanner", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ControllerOptions) ApplyTo(cfg *controller.ControllerConfig) error {
	var err error

	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst
	cfg.EnableRBAC = s.EnableRBAC

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
