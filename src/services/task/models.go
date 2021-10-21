package task

import "github.com/linshenqi/taskmate/src/base"

func (s *Service) dbListTasks(query *base.QueryTask) ([]*base.Task, int64, error) {
	tasks := []*base.Task{}
	var total int64 = 0

	countQuery := query.ToQuery(false)
	if err := countQuery.Find(&tasks).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQuery := query.ToQuery(true)
	if err := listQuery.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (s *Service) dbDeleteTasksByIDs(ids []string) error {
	if err := s.db.Delete(&base.Task{}, ids).Error; err != nil {
		return err
	}

	return nil
}
