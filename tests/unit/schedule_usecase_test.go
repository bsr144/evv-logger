package unit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/evv/schedule"
	apperrors "github.com/bsr144/evv-logger/pkg/errors"
)

type mockScheduleRepo struct {
	mock.Mock
}

func (m *mockScheduleRepo) CreateSchedule(ctx context.Context, s *entity.Schedule) (*entity.Schedule, error) {
	args := m.Called(ctx, s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Schedule), args.Error(1)
}

func (m *mockScheduleRepo) GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Schedule), args.Error(1)
}

func (m *mockScheduleRepo) GetAllSchedules(ctx context.Context) ([]*entity.Schedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Schedule), args.Error(1)
}

func (m *mockScheduleRepo) UpdateSchedule(ctx context.Context, s *entity.Schedule) (*entity.Schedule, error) {
	args := m.Called(ctx, s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Schedule), args.Error(1)
}

func (m *mockScheduleRepo) GetScheduleCounts(ctx context.Context, date string) (*entity.ScheduleCounts, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ScheduleCounts), args.Error(1)
}

type mockScheduleTaskRepo struct {
	mock.Mock
}

func (m *mockScheduleTaskRepo) CreateTask(ctx context.Context, t *entity.Task) (*entity.Task, error) {
	args := m.Called(ctx, t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockScheduleTaskRepo) GetTaskByID(ctx context.Context, id int) (*entity.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockScheduleTaskRepo) GetTasksByScheduleID(ctx context.Context, scheduleID int) ([]*entity.Task, error) {
	args := m.Called(ctx, scheduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *mockScheduleTaskRepo) UpdateTask(ctx context.Context, t *entity.Task) (*entity.Task, error) {
	args := m.Called(ctx, t)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func TestScheduleCreate_Success(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	schedRepo.On("CreateSchedule", ctx, mock.AnythingOfType("*entity.Schedule")).Return(
		&entity.Schedule{
			ID: 1, CaregiverName: "John", PatientName: "Jane",
			Date: "2026-03-08", StartTime: "09:00", EndTime: "17:00",
			Status: entity.ScheduleStatusUpcoming, Location: "123 Main St",
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}, nil,
	)

	resp, err := uc.Create(ctx, &schedule.CreateScheduleRequest{
		CaregiverName: "John", PatientName: "Jane",
		Date: "2026-03-08", StartTime: "09:00", EndTime: "17:00",
		Location: "123 Main St",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "John", resp.CaregiverName)
	assert.Equal(t, entity.ScheduleStatusUpcoming, resp.Status)
	schedRepo.AssertExpectations(t)
}

func TestScheduleCreate_MissingRequiredFields(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()

	tests := []struct {
		name string
		req  *schedule.CreateScheduleRequest
	}{
		{"missing caregiver", &schedule.CreateScheduleRequest{PatientName: "Jane", Date: "2026-03-08", StartTime: "09:00", EndTime: "17:00"}},
		{"missing patient", &schedule.CreateScheduleRequest{CaregiverName: "John", Date: "2026-03-08", StartTime: "09:00", EndTime: "17:00"}},
		{"missing date", &schedule.CreateScheduleRequest{CaregiverName: "John", PatientName: "Jane", StartTime: "09:00", EndTime: "17:00"}},
		{"missing start_time", &schedule.CreateScheduleRequest{CaregiverName: "John", PatientName: "Jane", Date: "2026-03-08", EndTime: "17:00"}},
		{"missing end_time", &schedule.CreateScheduleRequest{CaregiverName: "John", PatientName: "Jane", Date: "2026-03-08", StartTime: "09:00"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedRepo.On("CreateSchedule", ctx, mock.AnythingOfType("*entity.Schedule")).Return(
				&entity.Schedule{ID: 1, Status: entity.ScheduleStatusUpcoming, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil,
			).Once()

			resp, err := uc.Create(ctx, tt.req)
			require.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestScheduleGetByID_Success(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	schedRepo.On("GetScheduleByID", ctx, 1).Return(
		&entity.Schedule{ID: 1, CaregiverName: "John", PatientName: "Jane", Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now}, nil,
	)
	taskRepo.On("GetTasksByScheduleID", ctx, 1).Return([]*entity.Task{
		{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
	}, nil)

	resp, err := uc.GetByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 1, resp.Schedule.ID)
	assert.Len(t, resp.Tasks, 1)
	assert.Equal(t, "Check vitals", resp.Tasks[0].Description)
	schedRepo.AssertExpectations(t)
	taskRepo.AssertExpectations(t)
}

func TestScheduleGetByID_NotFound(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	schedRepo.On("GetScheduleByID", ctx, 999).Return(nil, apperrors.ErrScheduleNotFound)

	resp, err := uc.GetByID(ctx, 999)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
}

func TestScheduleGetAll_Success(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	schedRepo.On("GetAllSchedules", ctx).Return([]*entity.Schedule{
		{ID: 1, CaregiverName: "John", PatientName: "Jane", Date: "2026-03-08", Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now},
		{ID: 2, CaregiverName: "Bob", PatientName: "Alice", Date: "2026-03-09", Status: entity.ScheduleStatusCompleted, CreatedAt: now, UpdatedAt: now},
		{ID: 3, CaregiverName: "Eve", PatientName: "Mallory", Date: "2026-03-08", Status: entity.ScheduleStatusInProgress, CreatedAt: now, UpdatedAt: now},
	}, nil)

	resp, err := uc.GetAll(ctx, "")

	require.NoError(t, err)
	assert.Len(t, resp.Data, 3)
	assert.Equal(t, 3, resp.Stats.Total)
	assert.Equal(t, 1, resp.Stats.Upcoming)
	assert.Equal(t, 1, resp.Stats.InProgress)
	assert.Equal(t, 1, resp.Stats.Completed)
	assert.Equal(t, 0, resp.Stats.Missed)
}

func TestScheduleGetAll_WithDateFilter(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	schedRepo.On("GetAllSchedules", ctx).Return([]*entity.Schedule{
		{ID: 1, Date: "2026-03-08", Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Date: "2026-03-09", Status: entity.ScheduleStatusCompleted, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Date: "2026-03-08", Status: entity.ScheduleStatusInProgress, CreatedAt: now, UpdatedAt: now},
	}, nil)

	resp, err := uc.GetAll(ctx, "2026-03-08")

	require.NoError(t, err)
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 2, resp.Stats.Total)
	assert.Equal(t, 1, resp.Stats.Upcoming)
	assert.Equal(t, 1, resp.Stats.InProgress)
}

func TestScheduleGetAll_Empty(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	schedRepo.On("GetAllSchedules", ctx).Return([]*entity.Schedule{}, nil)

	resp, err := uc.GetAll(ctx, "")

	require.NoError(t, err)
	assert.Empty(t, resp.Data)
	assert.Equal(t, 0, resp.Stats.Total)
}

func TestScheduleUpdate_StatusTransition(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	status := entity.ScheduleStatusInProgress
	existing := &entity.Schedule{ID: 1, Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now}

	schedRepo.On("GetScheduleByID", ctx, 1).Return(existing, nil)
	schedRepo.On("UpdateSchedule", ctx, mock.AnythingOfType("*entity.Schedule")).Return(
		&entity.Schedule{ID: 1, Status: entity.ScheduleStatusInProgress, CreatedAt: now, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, &schedule.UpdateScheduleRequest{Status: &status})

	require.NoError(t, err)
	assert.Equal(t, entity.ScheduleStatusInProgress, resp.Status)
	schedRepo.AssertExpectations(t)
}

func TestScheduleUpdate_InvalidStatus(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	status := "invalid_status"
	existing := &entity.Schedule{ID: 1, Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now}

	schedRepo.On("GetScheduleByID", ctx, 1).Return(existing, nil)

	resp, err := uc.Update(ctx, 1, &schedule.UpdateScheduleRequest{Status: &status})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrInvalidStatus)
}

func TestScheduleUpdate_NotFound(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	status := entity.ScheduleStatusInProgress

	schedRepo.On("GetScheduleByID", ctx, 999).Return(nil, apperrors.ErrScheduleNotFound)

	resp, err := uc.Update(ctx, 999, &schedule.UpdateScheduleRequest{Status: &status})

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
}

func TestScheduleUpdate_ClockIn(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	clockIn := time.Now()
	lat := 40.7128
	lng := -74.0060
	existing := &entity.Schedule{ID: 1, Status: entity.ScheduleStatusUpcoming, CreatedAt: now, UpdatedAt: now}

	schedRepo.On("GetScheduleByID", ctx, 1).Return(existing, nil)
	schedRepo.On("UpdateSchedule", ctx, mock.AnythingOfType("*entity.Schedule")).Return(
		&entity.Schedule{ID: 1, Status: entity.ScheduleStatusUpcoming, ClockInTime: &clockIn, ClockInLat: &lat, ClockInLng: &lng, CreatedAt: now, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, &schedule.UpdateScheduleRequest{ClockInTime: &clockIn, ClockInLat: &lat, ClockInLng: &lng})

	require.NoError(t, err)
	assert.NotNil(t, resp.ClockInTime)
	assert.Equal(t, lat, *resp.ClockInLat)
	assert.Equal(t, lng, *resp.ClockInLng)
}

func TestScheduleUpdate_ClockOut(t *testing.T) {
	schedRepo := new(mockScheduleRepo)
	taskRepo := new(mockScheduleTaskRepo)
	uc := schedule.NewUsecase(schedRepo, taskRepo)

	ctx := context.Background()
	now := time.Now()
	clockOut := time.Now()
	lat := 40.7128
	lng := -74.0060
	existing := &entity.Schedule{ID: 1, Status: entity.ScheduleStatusInProgress, CreatedAt: now, UpdatedAt: now}

	schedRepo.On("GetScheduleByID", ctx, 1).Return(existing, nil)
	schedRepo.On("UpdateSchedule", ctx, mock.AnythingOfType("*entity.Schedule")).Return(
		&entity.Schedule{ID: 1, Status: entity.ScheduleStatusInProgress, ClockOutTime: &clockOut, ClockOutLat: &lat, ClockOutLng: &lng, CreatedAt: now, UpdatedAt: time.Now()}, nil,
	)

	resp, err := uc.Update(ctx, 1, &schedule.UpdateScheduleRequest{ClockOutTime: &clockOut, ClockOutLat: &lat, ClockOutLng: &lng})

	require.NoError(t, err)
	assert.NotNil(t, resp.ClockOutTime)
	assert.Equal(t, lat, *resp.ClockOutLat)
	assert.Equal(t, lng, *resp.ClockOutLng)
}
