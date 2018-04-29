package v1alpha1

func (r ScanResult) HasVulnerabilities() bool {
	for _, feature := range r.Features {
		if len(feature.Vulnerabilities) > 0 {
			return true
		}
	}
	return false
}

func (r ImageReviewResponse) HasVulnerabilities() bool {
	for _, img := range r.Images {
		if img.HasVulnerabilities() {
			return true
		}
	}
	return false
}
