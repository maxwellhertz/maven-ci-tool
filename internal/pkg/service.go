package pkg

import "io"

type PomUtils interface {
	// ConfigureJaCoCo configures jacoco-maven-plugin and maven-surefire-plugin
	// and returns the content of the updated POM file.
	ConfigurePlugins(pomFile io.Reader) ([]byte, error)
}

type MavenCmdRunner interface {
	// Run executes a Maven command.
	// cmd is the sub-command you want to run, such as "test".
	Run(cmd string, args ...string) error
}

type CliApp interface {
	Start() error	
}
