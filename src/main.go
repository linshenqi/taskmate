package main

import (
	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/services/instance"
	"github.com/linshenqi/taskmate/src/services/task"
)

func main() {

	app := sptty.GetApp()
	app.LoadConfFromFile()

	app.AddServices(sptty.Services{
		&task.Service{},
		&instance.Service{},
	})

	app.AddConfigs(sptty.Configs{})

	app.Sptting()
}
