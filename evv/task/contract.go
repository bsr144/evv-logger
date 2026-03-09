package task

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetTaskByID(ctx context.Context, id int) (*entity.Task, error)
	GetTasksByScheduleID(ctx context.Context, scheduleID int) ([]*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
}

type ScheduleChecker interface {
	GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error)
}
