package controller

import (
	"github.com/appscode/kutil/admission"
	workload "github.com/appscode/kutil/workload/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ admission.ResourceHandler = &ScannerController{}

func (c *ScannerController) OnCreate(obj runtime.Object) (runtime.Object, error) {
	w := obj.(*workload.Workload)
	w.Object = nil
	modObj, _, err := c.checkWorkload(w)
	return modObj, err
}

func (c *ScannerController) OnUpdate(oldObj, newObj runtime.Object) (runtime.Object, error) {
	w := newObj.(*workload.Workload)
	w.Object = nil
	modObj, _, err := c.checkWorkload(w)
	return modObj, err
}

func (c *ScannerController) OnDelete(obj runtime.Object) error {
	return nil
}
