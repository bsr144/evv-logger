package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

func (s *Store) GetScheduleCounts(_ context.Context, date string) (*entity.ScheduleCounts, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := &entity.ScheduleCounts{
		Total:            len(s.schedules),
		StatusCounts:     make(map[string]int),
		DateStatusCounts: make(map[string]int),
	}

	for status, count := range s.statusCounts {
		result.StatusCounts[status] = count
	}

	if dateCounts, ok := s.dateStatusCounts[date]; ok {
		for status, count := range dateCounts {
			result.DateStatusCounts[status] = count
		}
	}

	return result, nil
}
