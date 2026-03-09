package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

func (s *Store) GetScheduleByID(_ context.Context, id int) (*entity.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedule, ok := s.schedules[id]
	if !ok {
		return nil, pkgerrors.ErrScheduleNotFound
	}

	return deepCopySchedule(schedule), nil
}
