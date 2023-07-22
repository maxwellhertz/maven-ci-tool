package app_test

import (
	"testing"
	"time"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/app"
)

func TestGetLastestRelease(t *testing.T) {
	testCases := []struct {
		artifactId string
	}{
		{
			artifactId: "maven-surefire-plugin",
		},
		{
			artifactId: "jacoco-maven-plugin",
		},
	}

	api := setUp()

	for _, tc := range testCases {
		t.Run(tc.artifactId, func(t *testing.T) {
			version, err := api.GetLatestReleaseVersion(tc.artifactId)
			if err != nil {
				t.Errorf("expected no errors, but got: %v", err)
			}
			t.Logf("the version of the latest release of %v is %v", tc.artifactId, version)
		})
	}
}

func setUp() *app.MavenCetralRepoApi {
	return app.NewMavenCentralRepoApi(30 * time.Second).(*app.MavenCetralRepoApi)
}
