package memory

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

func (s *Store) CreateTask(_ context.Context, task *entity.Task) (*entity.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.taskNextID
	s.taskNextID++

	stored := *task
	s.tasks[stored.ID] = &stored

	result := stored
	return &result, nil
}
