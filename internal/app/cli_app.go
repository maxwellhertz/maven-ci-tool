package app

import (
	"fmt"
	"log"
	"os"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/pkg"
)

type DefaultCliApp struct {
	pomUtils       pkg.PomUtils
	mavenCmdRunner pkg.MavenCmdRunner
}

const (
	pomFilename       = "pom.xml"
	backupPomFilename = "pom_backup.xml"

	cleanCmd = "clean"
	testCmd  = "test"
)

func NewDefaultCliApp(pomUtils pkg.PomUtils, mavenCmdRunner pkg.MavenCmdRunner) pkg.CliApp {
	return &DefaultCliApp{pomUtils: pomUtils, mavenCmdRunner: mavenCmdRunner}
}

func (cli *DefaultCliApp) Start() error {
	if err := os.Rename(pomFilename, backupPomFilename); err != nil {
		return fmt.Errorf("can't rename %v to %v: %w", pomFilename, backupPomFilename, err)
	}
	defer func() {
		if err := os.Remove(pomFilename); err != nil {
			log.Fatalf("can't remove %v: %v", pomFilename, err)
		}
		if err := os.Rename(backupPomFilename, pomFilename); err != nil {
			log.Fatalf("can't rename %v back to %v: %v", backupPomFilename, pomFilename, err)
		}
	}()

	log.Println("Backing up the original " + pomFilename)
	backupPom, err := os.Open(backupPomFilename)
	if err != nil {
		return fmt.Errorf("can't open %v: %w", backupPomFilename, err)
	}
	defer backupPom.Close()

	log.Printf("Updating the %v to include necessary build plugins", pomFilename)
	updatedContent, err := cli.pomUtils.ConfigurePlugins(backupPom)
	if err != nil {
		return fmt.Errorf("can't update %v: %w", backupPomFilename, err)
	}

	pom, err := os.Create(pomFilename)
	if err != nil {
		return fmt.Errorf("can't create %v: %w", pomFilename, err)
	}
	defer pom.Close()
	if _, err := pom.Write(updatedContent); err != nil {
		return fmt.Errorf("can't save the updated POM to %v: %w", pomFilename, err)
	}

	log.Printf("Running Maven command `%v %v`", cleanCmd, testCmd)
	if err := cli.mavenCmdRunner.Run(cleanCmd, testCmd); err != nil {
		return fmt.Errorf("can't run Maven command: %w", err)
	}
	return nil
}
