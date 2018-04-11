package scanner

import (
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
)

type LayerType struct {
	Name       string
	Path       string
	ParentName string
	Format     string
	Features   []feature
	Headers    HeadersType
}

type HeadersType struct {
	Authorization string
}

type feature struct {
	Name            string              `json:"Name,omitempty"`
	NamespaceName   string              `json:"NamespaceName,omitempty"`
	Version         string              `json:"Version,omitempty"`
	Vulnerabilities []api.Vulnerability `json:"Vulnerabilities"`
	AddedBy         string              `json:"AddedBy,omitempty"`
}
