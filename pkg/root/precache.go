package root

import (
	"sync"

	"github.com/appscode/go/log"
	"github.com/soter/scanner/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InitSetter struct {
	ctl *controller.ScannerController
	//kubeClient kubernetes.Interface

	lock sync.RWMutex
}

func New(c *controller.ScannerController) *InitSetter {
	return &InitSetter{
		//kubeClient: kubeClient,
		ctl: c,
	}
}

func (s *InitSetter) Run() {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.scanRunningWorkloads()
}

func (s *InitSetter) scanRunningWorkloads() {
	ns, err := s.ctl.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Infoln(err)
	}

	for _, n := range ns.Items {
		s.scanDeployments(n.Name)
		s.scanReplicationControllers(n.Name)
		s.scanReplicaSets(n.Name)
		s.scanDaemonSets(n.Name)
		s.scanJobs(n.Name)
		s.scanStatefulSets(n.Name)
		//
		//pods, err := s.KubeClient.CoreV1().Pods(n.Name).List(metav1.ListOptions{})
		//if err != nil {
		//	log.Infof("in namespace(%s): %v", n.Name, err)
		//}
		//
		//
		//for _, pod := range pods.Items {
		//	for _, cont := range pod.Spec.Containers {
		//
		//	}
		//}
	}
}
