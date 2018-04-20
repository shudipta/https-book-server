package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/clair/api/v3/clairpb"
	api "github.com/soter/scanner/apis/scanner/v1alpha1"
	"google.golang.org/grpc"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func requestBearerToken(repo, userName, password string) (*http.Request, error) {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repo + ":pull&account=" + userName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if userName != "" {
		req.SetBasicAuth(userName, password)
	}

	return req, nil
}

func getBearerToken(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var token struct {
		Token string
	}

	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", token.Token), nil
}

func getVulnerabilities(res *clairpb.GetAncestryResponse) []api.Vulnerability {
	//return nil
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

func getFeaturs(res *clairpb.GetAncestryResponse) []api.Feature {
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

func clairClientSetup(clairAddress string) (clairpb.AncestryServiceClient, error) {
	conn, err := grpc.Dial(clairAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := clairpb.NewAncestryServiceClient(conn)
	return c, nil
}

func sendLayer(postAncestryRequest *clairpb.PostAncestryRequest, clairClient clairpb.AncestryServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err := clairClient.PostAncestry(ctx, postAncestryRequest)

	return err
}

func getLayer(
	repo string, clairClient clairpb.AncestryServiceClient) ([]api.Feature, []api.Vulnerability, error) {
	getAncestryRequest := &clairpb.GetAncestryRequest{
		AncestryName:        repo,
		WithFeatures:        true,
		WithVulnerabilities: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	resp, err := clairClient.GetAncestry(ctx, getAncestryRequest)
	if err != nil {
		return []api.Feature{}, []api.Vulnerability{}, err
	}

	return getFeaturs(resp), getVulnerabilities(resp), nil
}
