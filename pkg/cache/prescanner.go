package cache

import (
	"sync"

	"github.com/appscode/go/log"
	utilerrors "github.com/appscode/go/util/errors"
	"github.com/soter/scanner/pkg/controller"
)

type PreScanner struct {
	scanner *controller.ScannerController
	//kubeClient kubernetes.Interface

	lock sync.RWMutex
}

func New(c *controller.ScannerController) *PreScanner {
	return &PreScanner{
		//kubeClient: kubeClient,
		scanner: c,
	}
}

func (s *PreScanner) Run() {
	s.lock.RLock()
	defer s.lock.RUnlock()

	err := s.scanWorkloads()
	if err != nil {
		log.Errorln(err)
	}
}

func (s *PreScanner) scanWorkloads() error {
	var errs []error

	errs = append(errs, s.scanDeployments()...)
	errs = append(errs, s.scanReplicationControllers()...)
	errs = append(errs, s.scanReplicaSets()...)
	errs = append(errs, s.scanDaemonSets()...)
	errs = append(errs, s.scanJobs()...)
	errs = append(errs, s.scanCronJobs()...)
	errs = append(errs, s.scanStatefulSets()...)

	return utilerrors.NewAggregate(errs)
}
