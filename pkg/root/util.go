package root

import (
	"fmt"

	"github.com/appscode/go/log"
	workload "github.com/appscode/kubernetes-webhook-util/workload/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func logError(args ...interface{}) {
	log.Infoln(args)
}

func errorForGettingWorkloadFromNamespac(namespace string, err error) {
	logError(fmt.Sprintf("error in namespace(%s): %v", namespace, err))
}

func convertToWorkload(obj runtime.Object) *workload.Workload {
	w, err := workload.ConvertToWorkload(obj)
	if err != nil {
		logError("error in converting obj:", err)
	}

	return w
}
