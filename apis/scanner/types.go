package scanner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/api/core/v1"
)

type ImageReviewRequest struct {
	// +optional
	Image string

	// +optional
	ImagePullSecrets []core.LocalObjectReference
}

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
}

type ImageReviewResponse struct {
	// +optional
	Features []Feature

	// +optional
	Vulnerabilities []Vulnerability
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
