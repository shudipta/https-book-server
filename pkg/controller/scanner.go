package controller

import (
	workload "github.com/appscode/kutil/workload/v1"
)

func (c *ScannerController) checkWorkload(w *workload.Workload) (*workload.Workload, bool, error) {
	return w, false, nil
}
