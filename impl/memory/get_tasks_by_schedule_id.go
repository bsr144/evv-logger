package memory

import (
	"context"
	"sort"

	"github.com/bsr144/evv-logger/entity"
)

func (s *Store) GetTasksByScheduleID(_ context.Context, scheduleID int) ([]*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*entity.Task, 0)
	for _, task := range s.tasks {
		if task.ScheduleID == scheduleID {
			cp := *task
			result = append(result, &cp)
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })

	return result, nil
}
