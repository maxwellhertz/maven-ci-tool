package app_test

import (
	"testing"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/app"
)

func TestMavenCmdRunner(t *testing.T) {
	runner := app.NewDefaultMavenCmdRunner()
	_ = runner.Run("test")
}
