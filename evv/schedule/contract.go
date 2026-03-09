package schedule

import (
	"context"

	"github.com/bsr144/evv-logger/entity"
)

type ScheduleStatsResult struct {
	Total     int
	Missed    int
	Upcoming  int
	Completed int
}

type ScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error)
	GetAllSchedules(ctx context.Context) ([]*entity.Schedule, error)
	UpdateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleStats(ctx context.Context) (*ScheduleStatsResult, error)
}

type TaskRepository interface {
	GetTasksByScheduleID(ctx context.Context, scheduleID int) ([]*entity.Task, error)
}
