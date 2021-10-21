package base

import (
	"bytes"
	"fmt"
	"strings"
)

func ListAvailablePythonVersions() ([]string, error) {
	var output bytes.Buffer

	cmd, err := ShellExec("pyenv install --list", &output)
	if err != nil {
		return nil, err
	}

	_ = cmd.Wait()

	vals := strings.Split(output.String(), "\n")

	versions := []string{}
	for i := 1; i < len(vals); i++ {
		if vals[i] == "" {
			continue
		}

		versions = append(versions, strings.TrimSpace(vals[i]))
	}

	return versions, nil
}

func ListEnvs() ([]PythonEnv, error) {
	var output bytes.Buffer

	cmd, err := ShellExec("pyenv virtualenvs --bare --skip-aliases", &output)
	if err != nil {
		return nil, err
	}

	_ = cmd.Wait()

	envs := []PythonEnv{}
	vals := strings.Split(output.String(), "\n")

	for _, v := range vals {
		if v == "" {
			continue
		}

		kvs := strings.Split(strings.TrimSpace(v), "/envs/")
		envs = append(envs, PythonEnv{
			Version: kvs[0],
			Name:    kvs[1],
		})
	}

	return envs, nil
}

func CreateEnv(env *PythonEnv) error {
	if err := env.Validate(); err != nil {
		return err
	}

	var output bytes.Buffer

	cmd, err := ShellExec(fmt.Sprintf("pyenv install -s %s", env.Version), &output)
	if err != nil {
		return err
	}

	_ = cmd.Wait()

	cmd, err = ShellExec(fmt.Sprintf("pyenv virtualenv %s %s", env.Version, env.Name), &output)
	if err != nil {
		return err
	}

	_ = cmd.Wait()

	if err := ExecScript(env); err != nil {
		return err
	}

	return nil
}

func RemoveEnvs(names []string) error {
	var output bytes.Buffer
	for _, v := range names {

		cmd, err := ShellExec(fmt.Sprintf("pyenv virtualenv-delete -f %s ", v), &output)
		if err != nil {
			return err
		}

		_ = cmd.Wait()
	}

	return nil
}

func ExecScript(env *PythonEnv) error {
	if err := env.ValidateExecScript(); err != nil {
		return nil
	}

	var output bytes.Buffer

	cmd, err := ShellExec(fmt.Sprintf("eval \"$(pyenv init -)\" && eval \"$(pyenv virtualenv-init -)\" && pyenv activate %s && %s", env.Name, env.Script), &output)
	if err != nil {
		return err
	}

	_ = cmd.Wait()

	env.Msg = output.String()

	return nil
}
