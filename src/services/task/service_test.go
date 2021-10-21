package task

import (
	"fmt"
	"testing"

	"github.com/linshenqi/taskmate/src/base"
)

func TestListAvailablePythonVersions(t *testing.T) {
	versions, err := base.ListAvailablePythonVersions()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(versions)
}

func TestListEnvs(t *testing.T) {
	envs, err := base.ListEnvs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(envs)
}

func TestCreateEnv(t *testing.T) {
	err := base.CreateEnv(&base.PythonEnv{
		Version: "3.7.9",
		Name:    "test",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TestRemoveEnvs(t *testing.T) {
	err := base.RemoveEnvs([]string{
		"list",
		"test",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
