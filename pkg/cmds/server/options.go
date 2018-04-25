package server

import (
	"flag"

	"github.com/soter/scanner/pkg/controller"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ExtraOptions struct {
	ClairAddress string
	QPS          float64
	Burst        int
}

func NewExtraOptions() *ExtraOptions {
	return &ExtraOptions{
		ClairAddress: "http://clair.default.svc:6060",
		QPS:          100,
		Burst:        100,
	}
}

func (s *ExtraOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.ClairAddress, "clair-addr", s.ClairAddress, "The maximum QPS to the master from this client")

	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
}

func (s *ExtraOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("clair", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ExtraOptions) ApplyTo(cfg *controller.Config) error {
	var err error

	cfg.ClairAddress = s.ClairAddress
	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
