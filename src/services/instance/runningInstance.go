package instance

import "github.com/linshenqi/taskmate/src/base"

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
		delete(s.runningInstances, v)
	}
}
