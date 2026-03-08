package unit_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/impl/memory"
	apperrors "github.com/bsr144/evv-logger/pkg/errors"
)

func TestCreateSchedule(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	now := time.Now()

	s := &entity.Schedule{
		CaregiverName: "John", PatientName: "Jane",
		Date: "2026-03-08", StartTime: "09:00", EndTime: "17:00",
		Status: entity.ScheduleStatusUpcoming, Location: "123 Main St",
		CreatedAt: now, UpdatedAt: now,
	}

	created, err := store.CreateSchedule(ctx, s)

	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, "John", created.CaregiverName)
	assert.Equal(t, "Jane", created.PatientName)
	assert.Equal(t, entity.ScheduleStatusUpcoming, created.Status)
}

func TestCreateSchedule_AutoIncrementIDs(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	s1, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "A", PatientName: "B"})
	require.NoError(t, err)
	assert.Equal(t, 1, s1.ID)

	s2, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "C", PatientName: "D"})
	require.NoError(t, err)
	assert.Equal(t, 2, s2.ID)

	s3, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "E", PatientName: "F"})
	require.NoError(t, err)
	assert.Equal(t, 3, s3.ID)
}

func TestGetScheduleByID_Found(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "John", PatientName: "Jane", Status: entity.ScheduleStatusUpcoming})
	require.NoError(t, err)

	result, err := store.GetScheduleByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "John", result.CaregiverName)
}

func TestGetScheduleByID_NotFound(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	result, err := store.GetScheduleByID(ctx, 999)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
}

func TestGetAllSchedules(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "John", PatientName: "Jane"})
	require.NoError(t, err)
	_, err = store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "Bob", PatientName: "Alice"})
	require.NoError(t, err)

	results, err := store.GetAllSchedules(ctx)

	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestGetAllSchedules_Empty(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	results, err := store.GetAllSchedules(ctx)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestUpdateSchedule_Success(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "John", PatientName: "Jane", Status: entity.ScheduleStatusUpcoming})
	require.NoError(t, err)

	updated, err := store.UpdateSchedule(ctx, &entity.Schedule{ID: 1, CaregiverName: "John", PatientName: "Jane", Status: entity.ScheduleStatusInProgress})

	require.NoError(t, err)
	assert.Equal(t, entity.ScheduleStatusInProgress, updated.Status)
}

func TestUpdateSchedule_NotFound(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	result, err := store.UpdateSchedule(ctx, &entity.Schedule{ID: 999, Status: entity.ScheduleStatusInProgress})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrScheduleNotFound)
}

func TestCreateTask(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	now := time.Now()

	tk := &entity.Task{
		ScheduleID: 1, Description: "Check vitals",
		Status:    entity.TaskStatusPending,
		CreatedAt: now, UpdatedAt: now,
	}

	created, err := store.CreateTask(ctx, tk)

	require.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, 1, created.ScheduleID)
	assert.Equal(t, "Check vitals", created.Description)
	assert.Equal(t, entity.TaskStatusPending, created.Status)
}

func TestCreateTask_AutoIncrementIDs(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	t1, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Task A"})
	require.NoError(t, err)
	assert.Equal(t, 1, t1.ID)

	t2, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Task B"})
	require.NoError(t, err)
	assert.Equal(t, 2, t2.ID)

	t3, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 2, Description: "Task C"})
	require.NoError(t, err)
	assert.Equal(t, 3, t3.ID)
}

func TestGetTaskByID_Found(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending})
	require.NoError(t, err)

	result, err := store.GetTaskByID(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Check vitals", result.Description)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	result, err := store.GetTaskByID(ctx, 999)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrTaskNotFound)
}

func TestGetTasksByScheduleID(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Task A"})
	require.NoError(t, err)
	_, err = store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Task B"})
	require.NoError(t, err)
	_, err = store.CreateTask(ctx, &entity.Task{ScheduleID: 2, Description: "Task C"})
	require.NoError(t, err)

	results, err := store.GetTasksByScheduleID(ctx, 1)

	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Task A", results[0].Description)
	assert.Equal(t, "Task B", results[1].Description)
}

func TestGetTasksByScheduleID_Empty(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	results, err := store.GetTasksByScheduleID(ctx, 999)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestUpdateTask_Success(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending})
	require.NoError(t, err)

	updated, err := store.UpdateTask(ctx, &entity.Task{ID: 1, ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusCompleted, Notes: "Done"})

	require.NoError(t, err)
	assert.Equal(t, entity.TaskStatusCompleted, updated.Status)
	assert.Equal(t, "Done", updated.Notes)
}

func TestUpdateTask_NotFound(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	result, err := store.UpdateTask(ctx, &entity.Task{ID: 999, Status: entity.TaskStatusCompleted})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrTaskNotFound)
}

func TestDeepCopy_Schedule_MutationDoesNotAffectStore(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "John", PatientName: "Jane", Status: entity.ScheduleStatusUpcoming})
	require.NoError(t, err)

	result, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	result.CaregiverName = "MUTATED"
	result.Status = "MUTATED"

	original, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "John", original.CaregiverName)
	assert.Equal(t, entity.ScheduleStatusUpcoming, original.Status)
}

func TestDeepCopy_Task_MutationDoesNotAffectStore(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Check vitals", Status: entity.TaskStatusPending})
	require.NoError(t, err)

	result, err := store.GetTaskByID(ctx, 1)
	require.NoError(t, err)
	result.Description = "MUTATED"
	result.Status = "MUTATED"

	original, err := store.GetTaskByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Check vitals", original.Description)
	assert.Equal(t, entity.TaskStatusPending, original.Status)
}

func TestDeepCopy_GetAllSchedules_MutationDoesNotAffectStore(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "John", PatientName: "Jane"})
	require.NoError(t, err)

	results, err := store.GetAllSchedules(ctx)
	require.NoError(t, err)
	results[0].CaregiverName = "MUTATED"

	original, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "John", original.CaregiverName)
}

func TestDeepCopy_GetTasksByScheduleID_MutationDoesNotAffectStore(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: "Original"})
	require.NoError(t, err)

	results, err := store.GetTasksByScheduleID(ctx, 1)
	require.NoError(t, err)
	results[0].Description = "MUTATED"

	original, err := store.GetTaskByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Original", original.Description)
}

func TestConcurrentScheduleCreation(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	var wg sync.WaitGroup
	count := 100

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: fmt.Sprintf("Caregiver %d", n), PatientName: fmt.Sprintf("Patient %d", n)})
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	results, err := store.GetAllSchedules(ctx)
	require.NoError(t, err)
	assert.Len(t, results, count)

	ids := make(map[int]bool)
	for _, s := range results {
		assert.False(t, ids[s.ID], "duplicate ID found: %d", s.ID)
		ids[s.ID] = true
	}
}

func TestConcurrentTaskCreation(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	var wg sync.WaitGroup
	count := 100

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := store.CreateTask(ctx, &entity.Task{ScheduleID: 1, Description: fmt.Sprintf("Task %d", n)})
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	results, err := store.GetTasksByScheduleID(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, results, count)

	ids := make(map[int]bool)
	for _, tk := range results {
		assert.False(t, ids[tk.ID], "duplicate ID found: %d", tk.ID)
		ids[tk.ID] = true
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	var wg sync.WaitGroup

	_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: "Seed", PatientName: "Patient", Status: entity.ScheduleStatusUpcoming})
	require.NoError(t, err)

	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func(n int) {
			defer wg.Done()
			_, err := store.CreateSchedule(ctx, &entity.Schedule{CaregiverName: fmt.Sprintf("Writer %d", n), PatientName: "Patient"})
			require.NoError(t, err)
		}(i)
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_, err := store.GetAllSchedules(ctx)
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	results, err := store.GetAllSchedules(ctx)
	require.NoError(t, err)
	assert.Equal(t, 11, len(results))
}

func TestDeepCopy_SchedulePointerFields(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	clockIn := time.Now()
	lat := 40.7128
	lng := -74.0060

	_, err := store.CreateSchedule(ctx, &entity.Schedule{
		CaregiverName: "John", PatientName: "Jane",
		Status:      entity.ScheduleStatusInProgress,
		ClockInTime: &clockIn, ClockInLat: &lat, ClockInLng: &lng,
	})
	require.NoError(t, err)

	result, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	newLat := 0.0
	result.ClockInLat = &newLat

	original, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, lat, *original.ClockInLat)
}

func TestDeepCopy_ScheduleClockOutPointerFields(t *testing.T) {
	store := memory.New()
	ctx := context.Background()

	clockOut := time.Now()
	lat := 40.7128
	lng := -74.0060

	_, err := store.CreateSchedule(ctx, &entity.Schedule{
		CaregiverName: "John", PatientName: "Jane",
		Status:       entity.ScheduleStatusCompleted,
		ClockOutTime: &clockOut, ClockOutLat: &lat, ClockOutLng: &lng,
	})
	require.NoError(t, err)

	result, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	newTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	result.ClockOutTime = &newTime

	original, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, clockOut, *original.ClockOutTime)

	result2, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	newLat := 0.0
	result2.ClockOutLat = &newLat

	original2, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, lat, *original2.ClockOutLat)

	result3, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	newLng := 0.0
	result3.ClockOutLng = &newLng

	original3, err := store.GetScheduleByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, lng, *original3.ClockOutLng)
}
