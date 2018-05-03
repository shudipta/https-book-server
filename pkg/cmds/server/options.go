package server

import (
	"flag"

	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/soter/scanner/pkg/clair"
	"github.com/soter/scanner/pkg/controller"
	"github.com/soter/scanner/pkg/types"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type ExtraOptions struct {
	ClairAddress  string
	ClairCertDir  string
	QPS           float64
	Burst         int
	Severity      types.Severity
	FailurePolicy types.FailurePolicy
}

func NewExtraOptions() *ExtraOptions {
	return &ExtraOptions{
		ClairAddress:  "https://clairsvc.kube-system.svc:6060",
		ClairCertDir:  "/var/run/secrets/clair",
		QPS:           100,
		Burst:         100,
		Severity:      types.SeverityLow,
		FailurePolicy: types.FailurePolicyIgnore,
	}
}

func (s *ExtraOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.ClairAddress, "clair-addr", s.ClairAddress, "The address where clair is running")
	fs.StringVar(&s.ClairCertDir, "clair-cert-dir", s.ClairCertDir, "The directory where necessary certificates for clair api are stored")

	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
	fs.Var(&s.Severity, "highest-acceptable-severity",
		"Defines the highest acceptable level of vulnerability. If an image contains a vulnerable feature"+
			" with higher severity, it will be rejected by validating webhook. If an image contains vulnerabilities at this"+
			" level or lower, the image will be allowed to run.")
	fs.Var(&s.FailurePolicy, "failure-policy", "Define how errors from the docker registry are handled. Allowed values are Ignore or Fail.")
}

func (s *ExtraOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("scanner", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ExtraOptions) ApplyTo(cfg *controller.Config) error {
	var err error

	cfg.Severity = s.Severity
	cfg.FailurePolicy = s.FailurePolicy

	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.WorkloadClient, err = wcs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.Scanner, err = clair.NewScanner(cfg.ClientConfig, s.ClairAddress, s.ClairCertDir, s.Severity, s.FailurePolicy); err != nil {
		return err
	}

	return nil
}
