package instance

import (
	"fmt"
	"sync"

	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
	"gorm.io/gorm"
)

type Service struct {
	sptty.BaseService

	db          *gorm.DB
	serviceTask base.IServiceTask

	runningInstances map[string]*base.Instance
	mutex            sync.Mutex
}

func (s *Service) Init(app sptty.ISptty) error {
	s.db = app.Model().(*sptty.ModelService).DB()
	s.serviceTask = app.GetService(base.ServiceTask).(base.IServiceTask)

	s.runningInstances = map[string]*base.Instance{}
	s.mutex = sync.Mutex{}

	app.AddModel(&base.Instance{})

	app.AddRoute("POST", "/v1/instances", s.routePostInstances)
	app.AddRoute("GET", "/v1/instances", s.routeGetInstances)
	app.AddRoute("GET", "/v1/instances/{id:string}", s.routeGetInstancesByID)
	app.AddRoute("PUT", "/v1/instances-removal", s.routePutInstancesRemoval)

	return nil
}

func (s *Service) ServiceName() string {
	return base.ServiceInstance
}

func (s *Service) ListInstances(query *base.QueryInstance) ([]*base.Instance, int64, error) {
	instances, total, err := s.dbListInstances(query)
	if err != nil {
		return nil, 0, err
	}

	base.SerializeInstances(instances)

	return instances, total, nil
}

func (s *Service) GetInstanceByID(id string) (*base.Instance, error) {
	query := base.CreateQueryFromContext(&base.QueryInstance{
		QueryBase: base.QueryBase{
			IDs: []string{id},
		},
	}, s.db).(*base.QueryInstance)

	instances, _, err := s.ListInstances(query)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("Not Found")
	}

	targetInstance := instances[0].Serialize()

	runningInstance := s.tryGetRunningInstance(targetInstance)
	if runningInstance != nil {
		targetInstance.Msg = runningInstance.Output.String()
	}

	return targetInstance, nil
}

func (s *Service) tryGetRunningInstance(instance *base.Instance) *base.Instance {
	var runningInstance *base.Instance
	instances := s.getRunningInstances([]string{instance.ID})
	if len(instances) > 0 {
		runningInstance = instances[0]
	}

	return runningInstance
}

func (s *Service) CreateInstance(instance *base.Instance) error {
	tasks, _, _ := s.serviceTask.ListTasks(base.CreateQueryFromContext(&base.QueryTask{
		QueryBase: base.QueryBase{
			IDs: []string{instance.TaskID},
		},
	}, s.db).(*base.QueryTask))

	if len(tasks) == 0 {
		return fmt.Errorf("Task Not Found")
	}

	instance.Init()
	if err := s.db.Create(instance).Error; err != nil {
		return err
	}

	instance.Task = tasks[0]
	instance.Serialize()

	s.addRunningInstance(instance)
	go s.asyncRunInstance(instance)

	return nil
}

func (s *Service) DeleteInstances(ids []string) error {
	runningInstances := s.getRunningInstances(ids)
	for _, v := range runningInstances {
		if err := v.Stop(); err != nil {
			sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.Stop Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
		}
	}

	s.removeRunningInstances(ids)

	if err := s.dbDeleteInstancesByIDs(ids); err != nil {
		return err
	}

	return nil
}

func (s *Service) asyncRunInstance(instance *base.Instance) {
	if err := instance.Start(); err != nil {
		sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.Start Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
	}

	instance.Msg = instance.Output.String()
	if err := s.db.Updates(instance).Error; err != nil {
		sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.Updates Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
	}
}
