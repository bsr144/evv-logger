package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

func (s *Store) CreateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	schedule.ID = s.scheduleNextID
	s.scheduleNextID++

	stored := deepCopySchedule(schedule)
	s.schedules[stored.ID] = stored

	s.statusCounts[stored.Status]++
	if _, ok := s.dateStatusCounts[stored.Date]; !ok {
		s.dateStatusCounts[stored.Date] = make(map[string]int)
	}
	s.dateStatusCounts[stored.Date][stored.Status]++

	return deepCopySchedule(stored), nil
}
