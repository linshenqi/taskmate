package base

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	DefaultPageSize = 20
	HeaderTotal     = "Total"
)

const (
	StatusPending  = "pending"
	StatusRunning  = "running"
	StatusCreating = "creating"
	StatusSuccess  = "success"
	StatusFailed   = "failed"
)

func ShellExec(command string, output *bytes.Buffer) (*exec.Cmd, error) {
	cmd := exec.Command("/bin/bash", "-c", command)

	mw := io.MultiWriter(os.Stdout, output)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func CurrentFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	vals := strings.Split(f.Name(), ".")
	return vals[len(vals)-1]
}
