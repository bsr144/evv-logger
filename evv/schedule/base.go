package schedule

import "context"

type Usecase interface {
	Create(ctx context.Context, req *CreateScheduleRequest) (*ScheduleResponse, error)
	GetByID(ctx context.Context, id int) (*ScheduleDetailResponse, error)
	GetAll(ctx context.Context, dateFilter string) (*ScheduleListResponse, error)
	Update(ctx context.Context, id int, req *UpdateScheduleRequest) (*ScheduleResponse, error)
	GetStats(ctx context.Context) (*StatsResponse, error)
}

type scheduleUsecase struct {
	scheduleRepo ScheduleRepository
	taskRepo     TaskRepository
}

func NewUsecase(scheduleRepo ScheduleRepository, taskRepo TaskRepository) Usecase {
	return &scheduleUsecase{scheduleRepo: scheduleRepo, taskRepo: taskRepo}
}
