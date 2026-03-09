package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

func (s *Store) UpdateTask(_ context.Context, task *entity.Task) (*entity.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return nil, pkgerrors.ErrTaskNotFound
	}

	stored := *task
	s.tasks[stored.ID] = &stored

	result := stored
	return &result, nil
}
