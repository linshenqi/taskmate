package task

import (
	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
	"gorm.io/gorm"
)

type Service struct {
	sptty.BaseService

	db *gorm.DB
}

func (s *Service) Init(app sptty.ISptty) error {
	s.db = app.Model().(*sptty.ModelService).DB()

	app.AddModel(&base.Task{})

	app.AddRoute("GET", "/v1/python-versions", s.routeGetPythonVersions)
	app.AddRoute("GET", "/v1/envs", s.routeGetEnvs)
	app.AddRoute("POST", "/v1/envs", s.routePostEnvs)
	app.AddRoute("POST", "/v1/envs-script", s.routePostEnvsScript)
	app.AddRoute("PUT", "/v1/envs-removal", s.routePutEnvsRemoval)

	app.AddRoute("GET", "/v1/tasks", s.routeGetTasks)
	app.AddRoute("POST", "/v1/tasks", s.routePostTasks)
	app.AddRoute("PUT", "/v1/tasks-removal", s.routePutTasksRemoval)

	return nil
}

func (s *Service) ServiceName() string {
	return base.ServiceTask
}

func (s *Service) ListTasks(query *base.QueryTask) ([]*base.Task, int64, error) {
	tasks, total, err := s.dbListTasks(query)
	if err != nil {
		return nil, 0, err
	}

	base.SerializeTasks(tasks)

	return tasks, total, nil
}

func (s *Service) CreateTask(task *base.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	task.Init()
	if err := s.db.Create(task).Error; err != nil {
		return err
	}

	return nil
}
