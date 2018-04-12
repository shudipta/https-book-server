// +build !ignore_autogenerated

/*
Copyright 2018 The Soter Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	unsafe "unsafe"

	scanner "github.com/soter/scanner/apis/scanner"
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_Feature_To_scanner_Feature,
		Convert_scanner_Feature_To_v1alpha1_Feature,
		Convert_v1alpha1_ImageReview_To_scanner_ImageReview,
		Convert_scanner_ImageReview_To_v1alpha1_ImageReview,
		Convert_v1alpha1_ImageReviewRequest_To_scanner_ImageReviewRequest,
		Convert_scanner_ImageReviewRequest_To_v1alpha1_ImageReviewRequest,
		Convert_v1alpha1_ImageReviewResponse_To_scanner_ImageReviewResponse,
		Convert_scanner_ImageReviewResponse_To_v1alpha1_ImageReviewResponse,
		Convert_v1alpha1_Vulnerability_To_scanner_Vulnerability,
		Convert_scanner_Vulnerability_To_v1alpha1_Vulnerability,
	)
}

func autoConvert_v1alpha1_Feature_To_scanner_Feature(in *Feature, out *scanner.Feature, s conversion.Scope) error {
	out.Name = in.Name
	out.NamespaceName = in.NamespaceName
	out.Version = in.Version
	return nil
}

// Convert_v1alpha1_Feature_To_scanner_Feature is an autogenerated conversion function.
func Convert_v1alpha1_Feature_To_scanner_Feature(in *Feature, out *scanner.Feature, s conversion.Scope) error {
	return autoConvert_v1alpha1_Feature_To_scanner_Feature(in, out, s)
}

func autoConvert_scanner_Feature_To_v1alpha1_Feature(in *scanner.Feature, out *Feature, s conversion.Scope) error {
	out.Name = in.Name
	out.NamespaceName = in.NamespaceName
	out.Version = in.Version
	return nil
}

// Convert_scanner_Feature_To_v1alpha1_Feature is an autogenerated conversion function.
func Convert_scanner_Feature_To_v1alpha1_Feature(in *scanner.Feature, out *Feature, s conversion.Scope) error {
	return autoConvert_scanner_Feature_To_v1alpha1_Feature(in, out, s)
}

func autoConvert_v1alpha1_ImageReview_To_scanner_ImageReview(in *ImageReview, out *scanner.ImageReview, s conversion.Scope) error {
	out.Request = (*scanner.ImageReviewRequest)(unsafe.Pointer(in.Request))
	out.Response = (*scanner.ImageReviewResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_v1alpha1_ImageReview_To_scanner_ImageReview is an autogenerated conversion function.
func Convert_v1alpha1_ImageReview_To_scanner_ImageReview(in *ImageReview, out *scanner.ImageReview, s conversion.Scope) error {
	return autoConvert_v1alpha1_ImageReview_To_scanner_ImageReview(in, out, s)
}

func autoConvert_scanner_ImageReview_To_v1alpha1_ImageReview(in *scanner.ImageReview, out *ImageReview, s conversion.Scope) error {
	out.Request = (*ImageReviewRequest)(unsafe.Pointer(in.Request))
	out.Response = (*ImageReviewResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_scanner_ImageReview_To_v1alpha1_ImageReview is an autogenerated conversion function.
func Convert_scanner_ImageReview_To_v1alpha1_ImageReview(in *scanner.ImageReview, out *ImageReview, s conversion.Scope) error {
	return autoConvert_scanner_ImageReview_To_v1alpha1_ImageReview(in, out, s)
}

func autoConvert_v1alpha1_ImageReviewRequest_To_scanner_ImageReviewRequest(in *ImageReviewRequest, out *scanner.ImageReviewRequest, s conversion.Scope) error {
	out.Image = in.Image
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_v1alpha1_ImageReviewRequest_To_scanner_ImageReviewRequest is an autogenerated conversion function.
func Convert_v1alpha1_ImageReviewRequest_To_scanner_ImageReviewRequest(in *ImageReviewRequest, out *scanner.ImageReviewRequest, s conversion.Scope) error {
	return autoConvert_v1alpha1_ImageReviewRequest_To_scanner_ImageReviewRequest(in, out, s)
}

func autoConvert_scanner_ImageReviewRequest_To_v1alpha1_ImageReviewRequest(in *scanner.ImageReviewRequest, out *ImageReviewRequest, s conversion.Scope) error {
	out.Image = in.Image
	out.ImagePullSecrets = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	return nil
}

// Convert_scanner_ImageReviewRequest_To_v1alpha1_ImageReviewRequest is an autogenerated conversion function.
func Convert_scanner_ImageReviewRequest_To_v1alpha1_ImageReviewRequest(in *scanner.ImageReviewRequest, out *ImageReviewRequest, s conversion.Scope) error {
	return autoConvert_scanner_ImageReviewRequest_To_v1alpha1_ImageReviewRequest(in, out, s)
}

func autoConvert_v1alpha1_ImageReviewResponse_To_scanner_ImageReviewResponse(in *ImageReviewResponse, out *scanner.ImageReviewResponse, s conversion.Scope) error {
	out.Features = *(*[]scanner.Feature)(unsafe.Pointer(&in.Features))
	out.Vulnerabilities = *(*[]scanner.Vulnerability)(unsafe.Pointer(&in.Vulnerabilities))
	return nil
}

// Convert_v1alpha1_ImageReviewResponse_To_scanner_ImageReviewResponse is an autogenerated conversion function.
func Convert_v1alpha1_ImageReviewResponse_To_scanner_ImageReviewResponse(in *ImageReviewResponse, out *scanner.ImageReviewResponse, s conversion.Scope) error {
	return autoConvert_v1alpha1_ImageReviewResponse_To_scanner_ImageReviewResponse(in, out, s)
}

func autoConvert_scanner_ImageReviewResponse_To_v1alpha1_ImageReviewResponse(in *scanner.ImageReviewResponse, out *ImageReviewResponse, s conversion.Scope) error {
	out.Features = *(*[]Feature)(unsafe.Pointer(&in.Features))
	out.Vulnerabilities = *(*[]Vulnerability)(unsafe.Pointer(&in.Vulnerabilities))
	return nil
}

// Convert_scanner_ImageReviewResponse_To_v1alpha1_ImageReviewResponse is an autogenerated conversion function.
func Convert_scanner_ImageReviewResponse_To_v1alpha1_ImageReviewResponse(in *scanner.ImageReviewResponse, out *ImageReviewResponse, s conversion.Scope) error {
	return autoConvert_scanner_ImageReviewResponse_To_v1alpha1_ImageReviewResponse(in, out, s)
}

func autoConvert_v1alpha1_Vulnerability_To_scanner_Vulnerability(in *Vulnerability, out *scanner.Vulnerability, s conversion.Scope) error {
	out.Name = in.Name
	out.NamespaceName = in.NamespaceName
	out.Description = in.Description
	out.Link = in.Link
	out.Severity = in.Severity
	out.FixedBy = in.FixedBy
	out.FeatureName = in.FeatureName
	return nil
}

// Convert_v1alpha1_Vulnerability_To_scanner_Vulnerability is an autogenerated conversion function.
func Convert_v1alpha1_Vulnerability_To_scanner_Vulnerability(in *Vulnerability, out *scanner.Vulnerability, s conversion.Scope) error {
	return autoConvert_v1alpha1_Vulnerability_To_scanner_Vulnerability(in, out, s)
}

func autoConvert_scanner_Vulnerability_To_v1alpha1_Vulnerability(in *scanner.Vulnerability, out *Vulnerability, s conversion.Scope) error {
	out.Name = in.Name
	out.NamespaceName = in.NamespaceName
	out.Description = in.Description
	out.Link = in.Link
	out.Severity = in.Severity
	out.FixedBy = in.FixedBy
	out.FeatureName = in.FeatureName
	return nil
}

// Convert_scanner_Vulnerability_To_v1alpha1_Vulnerability is an autogenerated conversion function.
func Convert_scanner_Vulnerability_To_v1alpha1_Vulnerability(in *scanner.Vulnerability, out *Vulnerability, s conversion.Scope) error {
	return autoConvert_scanner_Vulnerability_To_v1alpha1_Vulnerability(in, out, s)
}
