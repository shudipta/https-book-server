package clair

import (
	"strings"

	"github.com/coreos/clair/api/v3/clairpb"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func getVulnerabilities(res *clairpb.GetAncestryResponse) []api.Vulnerability {
	var vuls []api.Vulnerability
	for _, feature := range res.Ancestry.Features {
		for _, vul := range feature.Vulnerabilities {
			vuls = append(vuls, api.Vulnerability{
				Name:          vul.Name,
				NamespaceName: vul.NamespaceName,
				Description:   vul.Description,
				Link:          vul.Link,
				Severity:      vul.Severity,
				FixedBy:       vul.FixedBy,
				FeatureName:   feature.Name,
			})
		}
	}

	return vuls
}

func getFeatures(res *clairpb.GetAncestryResponse) []api.Feature {
	var fs []api.Feature
	for _, feature := range res.Ancestry.Features {
		fs = append(fs, api.Feature{
			Name:          feature.Name,
			NamespaceName: feature.NamespaceName,
			Version:       feature.Version,
		})
	}

	return fs
}

func parseImageName(imageName string) (string, string, string, error) {
	repo, tag, digest, err := parsers.ParseImageName(imageName)
	if err != nil {
		return "", "", "", err
	}
	// the repo part should have registry url as prefix followed by a '/'
	// for example, if image name = "ubuntu" then
	//					repo = "docker.io/library/ubuntu", tag = "latest", digest = ""
	// 				for this image we need the repo = "library/ubuntu"
	//
	// 				if image name = "k8s.gcr.io/kubernetes-dashboard-amd64:v1.8.1" then
	//					repo = "k8s.gcr.io/kubernetes-dashboard-amd64", tag = "v1.8.1", digest = ""
	// 				for this image we need the repo = "kubernetes-dashboard-amd64"
	parts := strings.Split(repo, "/")
	repo = strings.Join(parts[1:], "/")

	return repo, tag, digest, err
}

func hashPart(digest string) string {
	if len(digest) < 7 {
		return ""
	}

	return digest[7:]
}
