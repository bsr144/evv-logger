package memory

import (
	"context"
	"sync"

	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/evv/schedule"
)

type Store struct {
	mu             sync.RWMutex
	schedules      map[int]*entity.Schedule
	tasks          map[int]*entity.Task
	scheduleNextID int
	taskNextID     int
}

func New() *Store {
	return &Store{
		schedules:      make(map[int]*entity.Schedule),
		tasks:          make(map[int]*entity.Task),
		scheduleNextID: 1,
		taskNextID:     1,
	}
}

var _ interface {
	CreateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleByID(_ context.Context, id int) (*entity.Schedule, error)
	GetAllSchedules(_ context.Context) ([]*entity.Schedule, error)
	UpdateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleStats(_ context.Context) (*schedule.ScheduleStatsResult, error)
	CreateTask(_ context.Context, task *entity.Task) (*entity.Task, error)
	GetTaskByID(_ context.Context, id int) (*entity.Task, error)
	GetTasksByScheduleID(_ context.Context, scheduleID int) ([]*entity.Task, error)
	UpdateTask(_ context.Context, task *entity.Task) (*entity.Task, error)
} = (*Store)(nil)
