package task

import (
	"context"
	"fmt"
	"time"

	"github.com/bsr144/evv-logger/entity"
	apperrors "github.com/bsr144/evv-logger/pkg/errors"
)

var validTaskStatuses = map[string]bool{
	entity.TaskStatusPending:      true,
	entity.TaskStatusCompleted:    true,
	entity.TaskStatusNotCompleted: true,
}

func (uc *taskUsecase) Create(ctx context.Context, scheduleID int, req *CreateTaskRequest) (*TaskResponse, error) {
	if req.Description == "" {
		return nil, fmt.Errorf("description is required: %w", apperrors.ErrValidation)
	}

	if _, err := uc.scheduleChecker.GetScheduleByID(ctx, scheduleID); err != nil {
		return nil, err
	}

	now := time.Now()
	task := &entity.Task{
		ScheduleID:  scheduleID,
		Description: req.Description,
		Status:      entity.TaskStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := uc.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return toTaskResponse(created), nil
}

func (uc *taskUsecase) GetByScheduleID(ctx context.Context, scheduleID int) ([]*TaskResponse, error) {
	if _, err := uc.scheduleChecker.GetScheduleByID(ctx, scheduleID); err != nil {
		return nil, err
	}

	tasks, err := uc.taskRepo.GetTasksByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	responses := make([]*TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		responses = append(responses, toTaskResponse(t))
	}

	return responses, nil
}

func (uc *taskUsecase) Update(ctx context.Context, scheduleID int, taskID int, req *UpdateTaskRequest) (*TaskResponse, error) {
	existing, err := uc.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if existing.ScheduleID != scheduleID {
		return nil, apperrors.ErrTaskNotFound
	}

	if req.Status != nil {
		if !validTaskStatuses[*req.Status] {
			return nil, fmt.Errorf("status '%s' is not valid: %w", *req.Status, apperrors.ErrInvalidStatus)
		}
		if *req.Status == entity.TaskStatusNotCompleted {
			if req.Notes == nil || *req.Notes == "" {
				return nil, fmt.Errorf("notes required when status is not_completed: %w", apperrors.ErrValidation)
			}
		}
		existing.Status = *req.Status
	}

	if req.Notes != nil {
		existing.Notes = *req.Notes
	}

	existing.UpdatedAt = time.Now()

	updated, err := uc.taskRepo.UpdateTask(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return toTaskResponse(updated), nil
}

func toTaskResponse(t *entity.Task) *TaskResponse {
	return &TaskResponse{
		ID:          t.ID,
		ScheduleID:  t.ScheduleID,
		Description: t.Description,
		Status:      t.Status,
		Notes:       t.Notes,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
