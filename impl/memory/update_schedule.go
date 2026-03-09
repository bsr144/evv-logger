package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

func (s *Store) UpdateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.schedules[schedule.ID]
	if !ok {
		return nil, pkgerrors.ErrScheduleNotFound
	}

	s.statusCounts[old.Status]--
	s.dateStatusCounts[old.Date][old.Status]--

	stored := deepCopySchedule(schedule)
	s.schedules[stored.ID] = stored

	s.statusCounts[stored.Status]++
	if _, ok := s.dateStatusCounts[stored.Date]; !ok {
		s.dateStatusCounts[stored.Date] = make(map[string]int)
	}
	s.dateStatusCounts[stored.Date][stored.Status]++

	return deepCopySchedule(stored), nil
}
