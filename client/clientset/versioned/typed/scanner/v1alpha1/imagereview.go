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
package v1alpha1

import (
	v1alpha1 "github.com/soter/scanner/apis/scanner/v1alpha1"
	rest "k8s.io/client-go/rest"
)

// ImageReviewsGetter has a method to return a ImageReviewInterface.
// A group's client should implement this interface.
type ImageReviewsGetter interface {
	ImageReviews(namespace string) ImageReviewInterface
}

// ImageReviewInterface has methods to work with ImageReview resources.
type ImageReviewInterface interface {
	Create(*v1alpha1.ImageReview) (*v1alpha1.ImageReview, error)
	ImageReviewExpansion
}

// imageReviews implements ImageReviewInterface
type imageReviews struct {
	client rest.Interface
	ns     string
}

// newImageReviews returns a ImageReviews
func newImageReviews(c *ScannerV1alpha1Client, namespace string) *imageReviews {
	return &imageReviews{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a imageReview and creates it.  Returns the server's representation of the imageReview, and an error, if there is any.
func (c *imageReviews) Create(imageReview *v1alpha1.ImageReview) (result *v1alpha1.ImageReview, err error) {
	result = &v1alpha1.ImageReview{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("imagereviews").
		Body(imageReview).
		Do().
		Into(result)
	return
}
