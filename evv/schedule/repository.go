package schedule

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

type ScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error)
	GetAllSchedules(ctx context.Context) ([]*entity.Schedule, error)
	UpdateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
}

type TaskRepository interface {
	GetTasksByScheduleID(ctx context.Context, scheduleID int) ([]*entity.Task, error)
}
