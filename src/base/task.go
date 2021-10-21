package base

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
	"gorm.io/gorm"
)

const (
	ServiceTask = "task"
)

const (
	ExecutorPython = "python"
	ExecutorBash   = "bash"
)

type IServiceTask interface {
	ListTasks(query *QueryTask) ([]*Task, int64, error)
}

type Task struct {
	sptty.SimpleModelBase

	Name     string `gorm:"size:32" json:"name"`
	Desc     string `json:"desc"`
	Executor string `gorm:"size:32" json:"executor"`

	Script string `json:"script"`

	// for python executor
	Env string `gorm:"size:64" json:"env"`
}

func (s *Task) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Name Is Required")
	}

	switch s.Executor {
	case ExecutorBash, ExecutorPython:
	default:
		return fmt.Errorf("Executor Error")
	}

	if s.Script == "" {
		return fmt.Errorf("Script Is Required")
	}

	return nil
}

func (s *Task) Serialize() *Task {
	s.SimpleModelBase.Serialize()

	return s
}

func (s *Task) Init() {
	s.SimpleModelBase.Init()
}

func SerializeTasks(tasks []*Task) []*Task {
	for k := range tasks {
		tasks[k] = tasks[k].Serialize()
	}

	return tasks
}

type QueryTask struct {
	QueryBase

	Name     string
	Executor string
	Env      string
}

func (s *QueryTask) fromCtx(ctx iris.Context) {
	s.QueryBase.fromCtx(ctx)

	s.Name = ctx.URLParam("name")
	s.Executor = ctx.URLParam("executor")
	s.Env = ctx.URLParam("env")
}

func (s *QueryTask) ToQuery(paging bool) *gorm.DB {
	q := s.QueryBase.ToQuery(paging)

	if s.Name != "" {
		q = q.Where("name like ?", fmt.Sprintf("%%%s%%", s.Name))
	}

	if s.Executor != "" {
		q = q.Where("executor = ?", s.Executor)
	}

	if s.Env != "" {
		q = q.Where("env = ?", s.Env)
	}

	q = q.Order("id desc")

	return q
}

type PythonEnv struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Script  string `json:"script"`
	Msg     string `json:"msg"`
}

func (s *PythonEnv) Validate() error {
	if s.Version == "" {
		return fmt.Errorf("Version Is Required")
	}

	if s.Name == "" {
		return fmt.Errorf("Name Is Required")
	}

	return nil
}

func (s *PythonEnv) ValidateExecScript() error {
	if s.Script == "" {
		return fmt.Errorf("Script Is Required")
	}

	if s.Name == "" {
		return fmt.Errorf("Name Is Required")
	}

	return nil
}
