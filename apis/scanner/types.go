package scanner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
)

type Vulnerability struct {
	Name          string
	NamespaceName string
	Description   string
	Link          string
	Severity      string
	//Metadata      map[string]interface{} `json:"Metadata,omitempty"`
	FixedBy string
	//FixedIn     []Feature `json:"FixedIn,omitempty"`
	FeatureName string
}

type Feature struct {
	Name          string
	NamespaceName string
	Version       string
	// +optional
	Vulnerabilities []Vulnerability
}

type ScanResult struct {
	Name string
	// +optional
	Features []Feature
}

type ImageReviewRequest struct {
	// +optional
	Image string

	// +optional
	ImagePullSecrets []core.LocalObjectReference
}

type ImageReviewResponse struct {
	Images []ScanResult `json:"images,omitempty"`
}

// +genclient
// +genclient:skipVerbs=list,update,patch,delete,deleteCollection,watch
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageReview describes a peer ping request/response.
type ImageReview struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// +optional
	Request *ImageReviewRequest
	// +optional
	Response *ImageReviewResponse
}
