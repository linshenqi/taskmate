package instance

import (
	"fmt"

	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
)

func (s *Service) getRunningInstances(ids []string) []*base.Instance {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	rt := []*base.Instance{}
	for _, v := range ids {
		ri, exist := s.runningInstances[v]
		if exist {
			rt = append(rt, ri)
		}
	}

	return rt
}

func (s *Service) addRunningInstance(instance *base.Instance) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.runningInstances[instance.ID] = instance
}

func (s *Service) removeRunningInstances(ids []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range ids {
		ri, exist := s.runningInstances[v]
		if exist {
			if err := ri.Stop(); err != nil {
				sptty.Log(sptty.ErrorLevel, fmt.Sprintf("%s.Stop Failed: %s", base.CurrentFuncName(), err.Error()), s.ServiceName())
			}

			delete(s.runningInstances, v)
		}
	}
}
