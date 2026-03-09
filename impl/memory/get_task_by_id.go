package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

func (s *Store) GetTaskByID(_ context.Context, id int) (*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, pkgerrors.ErrTaskNotFound
	}

	cp := *task
	return &cp, nil
}
