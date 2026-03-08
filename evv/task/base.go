package task

import "context"

type Usecase interface {
	Create(ctx context.Context, scheduleID int, req *CreateTaskRequest) (*TaskResponse, error)
	GetByScheduleID(ctx context.Context, scheduleID int) ([]*TaskResponse, error)
	Update(ctx context.Context, scheduleID int, taskID int, req *UpdateTaskRequest) (*TaskResponse, error)
}

type taskUsecase struct {
	taskRepo        TaskRepository
	scheduleChecker ScheduleChecker
}

func NewUsecase(taskRepo TaskRepository, scheduleChecker ScheduleChecker) Usecase {
	return &taskUsecase{taskRepo: taskRepo, scheduleChecker: scheduleChecker}
}
