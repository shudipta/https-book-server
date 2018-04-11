package scanner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/api/core/v1"
	"github.com/soter/scanner/pkg/scanner"
)

type ImageReviewRequest struct {
	// +optional
	Image string

	// +optional
	ImagePullSecrets []core.LocalObjectReference
}

type ImageReviewResponse struct {
	// +optional
	Features []scanner.Feature

	// +optional
	Vulnerabilities []scanner.Vulnerability
}

// +genclient
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageReview describes a peer ping request/response.
type ImageReview struct {
	metav1.TypeMeta

	// +optional
	Request *ImageReviewRequest

	// +optional
	Response *ImageReviewResponse
}
