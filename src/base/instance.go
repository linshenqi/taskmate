package base

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
	"gorm.io/gorm"
)

const (
	ServiceInstance = "instance"
)

// const (
// 	InstanceTypeOneTime = "onetime"
// 	InstanceTypeLoop    = "loop"
// )

type InstanceContext struct {
	Cmd    *exec.Cmd    `gorm:"-" json:"-"`
	Done   chan bool    `gorm:"-" json:"-"`
	Output bytes.Buffer `gorm:"-" json:"-"`
}

type Instance struct {
	sptty.SimpleModelBase
	InstanceContext

	TaskID string `gorm:"size:32" json:"task_id,omitempty"`
	Task   *Task  `gorm:"foreignkey:TaskID" json:"task"`

	Params string `json:"params"`

	// Type string `gorm:"size:32" json:"type"`

	// for loop type
	// Duration time.Duration `json:"duration"`

	Status   string `gorm:"size:32" json:"status"`
	Msg      string `json:"msg"`
	WorkerID string `gorm:"size:32" json:"-"`
}

func SerializeInstances(instances []*Instance) []*Instance {
	for k := range instances {
		instances[k] = instances[k].Serialize()
	}

	return instances
}

func (s *Instance) Serialize() *Instance {
	s.SimpleModelBase.Serialize()

	s.TaskID = ""
	return s
}

func (s *Instance) Init() {
	s.SimpleModelBase.Init()

	s.Status = StatusPending
}

func (s *Instance) generateFile() (string, error) {
	if s.ID == "" {
		return "", fmt.Errorf("ID Is Required")
	}

	fileName := fmt.Sprintf("%s.py", s.ID)
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(s.Task.Script)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *Instance) Start() error {
	var err error

	fileName, err := s.generateFile()
	if err != nil {
		return err
	}

	defer os.Remove(fileName)

	cmd := fmt.Sprintf("eval \"$(pyenv init -)\" && eval \"$(pyenv virtualenv-init -)\" && pyenv activate %s && python %s %s", s.Task.Env, fileName, s.Params)
	s.Cmd, err = ShellExec(cmd, &s.Output)
	if err != nil {
		s.Status = StatusFailed
		return err
	}

	if err := s.Cmd.Wait(); err != nil {
		s.Status = StatusFailed
		return err
	}

	s.Status = StatusSuccess

	return nil
}

func (s *Instance) Stop() error {
	if s.Cmd == nil || s.Cmd.Process == nil {
		return nil
	}

	if err := s.Cmd.Process.Kill(); err != nil {
		return err
	}

	return nil
}

type QueryInstance struct {
	QueryBase

	TaskID string
	Status string
}

func (s *QueryInstance) fromCtx(ctx iris.Context) {
	s.QueryBase.fromCtx(ctx)

	s.TaskID = ctx.URLParam("task_id")
	s.Status = ctx.URLParam("status")
}

func (s *QueryInstance) ToQuery(paging bool) *gorm.DB {
	q := s.QueryBase.ToQuery(paging)

	if s.TaskID != "" {
		q = q.Where("task_id = ?", s.TaskID)
	}

	if s.Status != "" {
		q = q.Where("status = ?", s.Status)
	}

	q = q.Preload("Task").Order("id desc")

	return q
}
