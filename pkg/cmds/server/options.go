package server

import (
	"flag"

	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/soter/scanner/pkg/controller"
	"github.com/soter/scanner/pkg/types"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ExtraOptions struct {
	ClairAddress    string
	ClairApiCertDir string
	QPS             float64
	Burst           int
	FailurePolicy   types.FailurePolicy
}

func NewExtraOptions() *ExtraOptions {
	return &ExtraOptions{
		ClairAddress:    "http://clairsvc.default.svc:6060",
		ClairApiCertDir: "/var/clairapi-client-cert/",
		QPS:             100,
		Burst:           100,
		FailurePolicy:   types.FailurePolicyIgnore,
	}
}

func (s *ExtraOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.ClairAddress, "clair-addr", s.ClairAddress, "The address where clair is running")
	fs.StringVar(&s.ClairApiCertDir, "api-cert-dir", s.ClairApiCertDir, "The directory where necessary certificates for clair api are stored")

	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
	fs.Var(&s.FailurePolicy, "failure-policy", "Define how errors from the docker registry are handled. Allowed values are Ignore or Fail.")
}

func (s *ExtraOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("clair", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ExtraOptions) ApplyTo(cfg *controller.Config) error {
	var err error

	cfg.ClairAddress = s.ClairAddress
	cfg.ClairApiCertDir = s.ClairApiCertDir
	cfg.FailurePolicy = s.FailurePolicy

	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.WorkloadClient, err = wcs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
