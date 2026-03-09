package memory

import (
	"context"
	"sort"

	"github.com/bsr144/evv-logger/entity"
)

func (s *Store) GetAllSchedules(_ context.Context) ([]*entity.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*entity.Schedule, 0, len(s.schedules))
	for _, schedule := range s.schedules {
		result = append(result, deepCopySchedule(schedule))
	}

	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })

	return result, nil
}
