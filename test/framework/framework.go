package framework

import (
	"github.com/appscode/go/crypto/rand"
	"k8s.io/client-go/kubernetes"
	ka "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

type Framework struct {
	KubeClient     kubernetes.Interface
	KAClient       ka.Interface
	namespace      string
	WebhookEnabled bool
}

func New(kubeClient kubernetes.Interface, kaClient ka.Interface, webhookEnabled bool) *Framework {
	return &Framework{
		KubeClient:     kubeClient,
		KAClient:       kaClient,
		namespace:      rand.WithUniqSuffix("scanner-e2e"),
		WebhookEnabled: webhookEnabled,
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("test-scanner"),
	}
}

type Invocation struct {
	*Framework
	app string
}

func (f *Invocation) App() string {
	return f.app
}
