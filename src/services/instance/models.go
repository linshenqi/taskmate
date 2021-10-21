package instance

import "github.com/linshenqi/taskmate/src/base"

func (s *Service) dbListInstances(query *base.QueryInstance) ([]*base.Instance, int64, error) {
	instances := []*base.Instance{}
	var total int64 = 0

	countQuery := query.ToQuery(false)
	if err := countQuery.Find(&instances).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQuery := query.ToQuery(true)
	if err := listQuery.Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	return instances, total, nil
}

func (s *Service) dbDeleteInstancesByIDs(ids []string) error {
	if err := s.db.Delete(&base.Instance{}, ids).Error; err != nil {
		return err
	}

	return nil
}
