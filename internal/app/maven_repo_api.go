package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MavenCetralRepoApi struct {
	httpCli *http.Client
}

func NewMavenCentralRepoApi(timeout time.Duration) MavenRepositoryApi {
	return &MavenCetralRepoApi{
		httpCli: &http.Client{
			Timeout: timeout,
		},
	}
}

const (
	mavenCentralRepoUrl = "https://central.sonatype.com/api/internal/browse/components"
)

type mavenRepoRequest struct {
	SearchTerm string `json:"searchTerm"`
	Size       int    `json:"size"`
}

type mavenRepoResponse struct {
	Total      int `json:"totalResultCount"`
	Components []struct {
		LatestRelease struct {
			Version string `json:"version"`
		} `json:"latestVersionInfo"`
	} `json:"components"`
}

func (api *MavenCetralRepoApi) GetLatestReleaseVersion(artifactId string) (string, error) {
	req := mavenRepoRequest{
		SearchTerm: artifactId,
		Size:       1,
	}
	reqPayload, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := api.httpCli.Post(mavenCentralRepoUrl, "application/json", bytes.NewBuffer(reqPayload))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	var respPayload mavenRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		return "", err
	}
	if respPayload.Total == 0 {
		return "", fmt.Errorf("%v does not exist in the Maven cetral repository", artifactId)
	}
	return respPayload.Components[0].LatestRelease.Version, nil
}
