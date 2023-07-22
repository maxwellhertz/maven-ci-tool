package main

import (
	"log"
	"time"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/app"
)

func main() {
	pomUtils := app.NewDefaultPomUitls(app.NewMavenCentralRepoApi(30 * time.Second))
	cmdRunner := app.NewDefaultMavenCmdRunner()
	cli := app.NewDefaultCliApp(pomUtils, cmdRunner)
	if err := cli.Start(); err != nil {
		log.Fatalln(err)
	}
}
