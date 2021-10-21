package instance

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
)

func TestInstance(t *testing.T) {
	script, err := ioutil.ReadFile("../../../test/test.py")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	instance := base.Instance{
		SimpleModelBase: sptty.SimpleModelBase{
			ID: sptty.GenerateUID(),
		},

		Task: &base.Task{
			Executor: base.ExecutorPython,
			Env:      "qlib",
			Script:   string(script),
		},

		Params: "1234",
	}

	if err := instance.Start(); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(instance.Status)
}
