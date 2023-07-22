package app

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/pkg"
)

type DefaultMavenCommandRunner struct {
}

const (
	mvnExecName         = "mvn"
	mvnwExecName        = "./mvnw"
	mvnwWindowsExecName = "./mvnw.cmd"
)

func NewDefaultMavenCmdRunner() pkg.MavenCmdRunner {
	return &DefaultMavenCommandRunner{}
}

func (r *DefaultMavenCommandRunner) Run(mavenSubCmd string, args ...string) error {
	var mavenExec string
	switch {
	// Use the Maven wrapper if it exists.
	case runtime.GOOS == "windows" && fileExists(mvnwWindowsExecName):
		mavenExec = mvnwWindowsExecName
	case runtime.GOOS != "windows" && fileExists(mvnwExecName):
		mavenExec = mvnwExecName
	default:
		mavenExec = mvnExecName
	}

	cmd := exec.Command(mavenExec, append([]string{mavenSubCmd}, args...)...)
	// Connect the command's output directly to os.Stdout
	cmd.Stdout = os.Stdout
	// Connect the command's error output directly to os.Stderr
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
