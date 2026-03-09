package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

func (s *Store) UpdateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.schedules[schedule.ID]; !ok {
		return nil, pkgerrors.ErrScheduleNotFound
	}

	stored := deepCopySchedule(schedule)
	s.schedules[stored.ID] = stored

	return deepCopySchedule(stored), nil
}
