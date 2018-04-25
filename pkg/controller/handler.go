package controller

import (
	"github.com/appscode/kubernetes-webhook-util/admission"
	workload "github.com/appscode/kubernetes-webhook-util/apis/workload/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ admission.ResourceHandler = &Controller{}

func (c *Controller) OnCreate(obj runtime.Object) (runtime.Object, error) {
	w := obj.(*workload.Workload)
	w.Object = nil
	modObj, _, err := c.CheckWorkload(w)
	return modObj, err
}

func (c *Controller) OnUpdate(oldObj, newObj runtime.Object) (runtime.Object, error) {
	w := newObj.(*workload.Workload)
	w.Object = nil
	modObj, _, err := c.CheckWorkload(w)
	return modObj, err
}

func (c *Controller) OnDelete(obj runtime.Object) error {
	return nil
}
