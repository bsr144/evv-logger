package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/evv/task"
	apperrors "github.com/bsr144/evv-logger/pkg/errors"
)

type mockTaskRepo struct {
	mock.Mock
}

func (m *mockTaskRepo) CreateTask(ctx context.Context, t *entity.Task) (*entity.Task, error) {
	args := m.Called(ctx, t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockTaskRepo) GetTaskByID(ctx context.Context, id int) (*entity.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockTaskRepo) GetTasksByScheduleID(ctx context.Context, scheduleID int) ([]*entity.Task, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *mockTaskRepo) UpdateTask(ctx context.Context, t *entity.Task) (*entity.Task, error) {
	args := m.Called(ctx, t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

type mockScheduleChecker struct {
	mock.Mock
}

func (m *mockScheduleChecker) GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Schedule), args.Error(1)
}

func TestTaskCreate_Success(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	schedChecker.On("GetScheduleByID", ctx, 1).Return(&entity.Schedule{ID: 1}, nil)
	taskRepo.On("CreateTask", ctx, mock.AnythingOfType("*entity.Task")).Return(
		&entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Create(ctx, 1, &task.CreateTaskRequest{Description: "Check vitals"})

	require.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, 1, resp.ScheduleID)
	assert.Equal(t, "Check vitals", resp.Description)
	assert.Equal(t, entity.TaskStatusPending, resp.Status)
	schedChecker.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestTaskCreate_ScheduleNotFound(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	schedChecker.On("GetScheduleByID", ctx, 999).Return(nil, apperrors.ErrScheduleNotFound)

	resp, err := uc.Create(ctx, 999, &task.CreateTaskRequest{Description: "Check vitals"})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
	schedChecker.AssertExpectations(t)
}

func TestTaskCreate_EmptyDescription(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()

	resp, err := uc.Create(ctx, 1, &task.CreateTaskRequest{Description: ""})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestTaskGetByScheduleID_Success(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	now := time.Now()
	schedChecker.On("GetScheduleByID", ctx, 1).Return(&entity.Schedule{ID: 1}, nil)
	taskRepo.On("GetTasksByScheduleID", ctx, 1).Return([]*entity.Task{
		{ID: 1, ScheduleID: 1, Description: "Task A", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ID: 2, ScheduleID: 1, Description: "Task B", Status: entity.TaskStatusCompleted, CreatedAt: now, UpdatedAt: now},
	}, nil)

	resp, err := uc.GetByScheduleID(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Task A", resp[0].Description)
	assert.Equal(t, "Task B", resp[1].Description)
	schedChecker.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestTaskGetByScheduleID_ScheduleNotFound(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	schedChecker.On("GetScheduleByID", ctx, 999).Return(nil, apperrors.ErrScheduleNotFound)

	resp, err := uc.GetByScheduleID(ctx, 999)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
}

func TestTaskGetByScheduleID_Empty(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	schedChecker.On("GetScheduleByID", ctx, 1).Return(&entity.Schedule{ID: 1}, nil)
	taskRepo.On("GetTasksByScheduleID", ctx, 1).Return([]*entity.Task{}, nil)

	resp, err := uc.GetByScheduleID(ctx, 1)

	require.NoError(t, err)
	assert.Empty(t, resp)
}

func TestTaskUpdate_StatusToCompleted(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	status := entity.TaskStatusCompleted
	existing := &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	taskRepo.On("GetTaskByID", ctx, 1).Return(existing, nil)
	taskRepo.On("UpdateTask", ctx, mock.AnythingOfType("*entity.Task")).Return(
		&entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusCompleted, CreatedAt: existing.CreatedAt, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, 1, &task.UpdateTaskRequest{Status: &status})

	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusCompleted, resp.Status)
	taskRepo.AssertExpectations(t)
}

func TestTaskUpdate_StatusToNotCompleted_WithNotes(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	status := entity.TaskStatusNotCompleted
	notes := "Patient refused"
	existing := &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	taskRepo.On("GetTaskByID", ctx, 1).Return(existing, nil)
	taskRepo.On("UpdateTask", ctx, mock.AnythingOfType("*entity.Task")).Return(
		&entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusNotCompleted, Notes: "Patient refused", CreatedAt: existing.CreatedAt, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, 1, &task.UpdateTaskRequest{Status: &status, Notes: &notes})

	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusNotCompleted, resp.Status)
	assert.Equal(t, "Patient refused", resp.Notes)
}

func TestTaskUpdate_StatusToNotCompleted_WithoutNotes(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	status := entity.TaskStatusNotCompleted
	existing := &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	taskRepo.On("GetTaskByID", ctx, 1).Return(existing, nil)

	resp, err := uc.Update(ctx, 1, 1, &task.UpdateTaskRequest{Status: &status})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestTaskUpdate_InvalidStatus(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	status := "invalid_status"
	existing := &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	taskRepo.On("GetTaskByID", ctx, 1).Return(existing, nil)

	resp, err := uc.Update(ctx, 1, 1, &task.UpdateTaskRequest{Status: &status})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrInvalidStatus)
}

func TestTaskUpdate_TaskNotFound(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	status := entity.TaskStatusCompleted

	taskRepo.On("GetTaskByID", ctx, 999).Return(nil, apperrors.ErrTaskNotFound)

	resp, err := uc.Update(ctx, 1, 999, &task.UpdateTaskRequest{Status: &status})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrTaskNotFound)
}

func TestTaskUpdate_NotesOnlyWithoutStatusChange(t *testing.T) {
	taskRepo := new(mockTaskRepo)
	schedChecker := new(mockScheduleChecker)
	uc := task.NewUsecase(taskRepo, schedChecker)

	ctx := context.Background()
	notes := "Additional observations"
	existing := &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	taskRepo.On("GetTaskByID", ctx, 1).Return(existing, nil)
	taskRepo.On("UpdateTask", ctx, mock.AnythingOfType("*entity.Task")).Return(
		&entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, Notes: "Additional observations", CreatedAt: existing.CreatedAt, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, 1, &task.UpdateTaskRequest{Notes: &notes})

	require.NoError(t, err)
	assert.Equal(t, "Additional observations", resp.Notes)
	assert.Equal(t, entity.TaskStatusPending, resp.Status)
}
