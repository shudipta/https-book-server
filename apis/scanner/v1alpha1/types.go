package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Vulnerability represents vulnerability entity returned by Clair
type Vulnerability struct {
	Name          string `json:"Name,omitempty"`
	NamespaceName string `json:"NamespaceName,omitempty"`
	Description   string `json:"Description,omitempty"`
	Link          string `json:"Link,omitempty"`
	Severity      string `json:"Severity,omitempty"`
	//Metadata      map[string]interface{} `json:"Metadata,omitempty"`
	FixedBy string `json:"FixedBy,omitempty"`
	//FixedIn     []Feature `json:"FixedIn,omitempty"`
	FeatureName string `json:"featureName,omitempty"`
}

type Feature struct {
	Name          string `json:"Name,omitempty"`
	NamespaceName string `json:"NamespaceName,omitempty"`
	Version       string `json:"Version,omitempty"`
	// +optional
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
}

type ScanResult struct {
	Name string `json:"name,omitempty"`
	// +optional
	Features []Feature `json:"features,omitempty"`
}

type ImageReviewResponse struct {
	Images []ScanResult `json:"images,omitempty"`
}

const (
	ResourceKindImageReview     = "ImageReview"
	ResourcePluralImageReview   = "imagereviews"
	ResourceSingularImageReview = "imagereview"
)

// +genclient
// +genclient:skipVerbs=list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageReview describes a peer ping request/response.
type ImageReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Response *ImageReviewResponse `json:"response,omitempty"`
}
