package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/bsr144/evv-logger/entity"
	apperrors "github.com/bsr144/evv-logger/pkg/errors"
)

var validScheduleStatuses = map[string]bool{
	entity.ScheduleStatusUpcoming:   true,
	entity.ScheduleStatusInProgress: true,
	entity.ScheduleStatusCompleted:  true,
	entity.ScheduleStatusMissed:     true,
}

func (uc *scheduleUsecase) Create(ctx context.Context, req *CreateScheduleRequest) (*ScheduleResponse, error) {
	now := time.Now()
	schedule := &entity.Schedule{
		CaregiverName: req.CaregiverName,
		PatientName:   req.PatientName,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Location:      req.Location,
		Status:        entity.ScheduleStatusUpcoming,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	created, err := uc.scheduleRepo.CreateSchedule(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return toScheduleResponse(created), nil
}

func (uc *scheduleUsecase) GetByID(ctx context.Context, id int) (*ScheduleDetailResponse, error) {
	s, err := uc.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tasks, err := uc.taskRepo.GetTasksByScheduleID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for schedule: %w", err)
	}

	taskSummaries := make([]*TaskSummary, 0, len(tasks))
	for _, t := range tasks {
		taskSummaries = append(taskSummaries, &TaskSummary{
			ID:          t.ID,
			ScheduleID:  t.ScheduleID,
			Description: t.Description,
			Status:      t.Status,
			Notes:       t.Notes,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return &ScheduleDetailResponse{
		Schedule: toScheduleResponse(s),
		Tasks:    taskSummaries,
	}, nil
}

func (uc *scheduleUsecase) GetAll(ctx context.Context, dateFilter string) (*ScheduleListResponse, error) {
	schedules, err := uc.scheduleRepo.GetAllSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	var filtered []*entity.Schedule
	if dateFilter != "" {
		for _, s := range schedules {
			if s.Date == dateFilter {
				filtered = append(filtered, s)
			}
		}
	} else {
		filtered = schedules
	}

	data := make([]*ScheduleResponse, 0, len(filtered))
	stats := ScheduleStats{Total: len(filtered)}

	for _, s := range filtered {
		data = append(data, toScheduleResponse(s))
		switch s.Status {
		case entity.ScheduleStatusUpcoming:
			stats.Upcoming++
		case entity.ScheduleStatusInProgress:
			stats.InProgress++
		case entity.ScheduleStatusCompleted:
			stats.Completed++
		case entity.ScheduleStatusMissed:
			stats.Missed++
		}
	}

	return &ScheduleListResponse{Data: data, Stats: stats}, nil
}

func (uc *scheduleUsecase) Update(ctx context.Context, id int, req *UpdateScheduleRequest) (*ScheduleResponse, error) {
	existing, err := uc.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		if !validScheduleStatuses[*req.Status] {
			return nil, fmt.Errorf("status '%s' is not valid: %w", *req.Status, apperrors.ErrInvalidStatus)
		}
		existing.Status = *req.Status
	}

	if req.ClockInTime != nil {
		existing.ClockInTime = req.ClockInTime
	}
	if req.ClockInLat != nil {
		existing.ClockInLat = req.ClockInLat
	}
	if req.ClockInLng != nil {
		existing.ClockInLng = req.ClockInLng
	}
	if req.ClockOutTime != nil {
		existing.ClockOutTime = req.ClockOutTime
	}
	if req.ClockOutLat != nil {
		existing.ClockOutLat = req.ClockOutLat
	}
	if req.ClockOutLng != nil {
		existing.ClockOutLng = req.ClockOutLng
	}

	existing.UpdatedAt = time.Now()

	updated, err := uc.scheduleRepo.UpdateSchedule(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return toScheduleResponse(updated), nil
}

func toScheduleResponse(s *entity.Schedule) *ScheduleResponse {
	return &ScheduleResponse{
		ID:            s.ID,
		CaregiverName: s.CaregiverName,
		PatientName:   s.PatientName,
		Date:          s.Date,
		StartTime:     s.StartTime,
		EndTime:       s.EndTime,
		Status:        s.Status,
		ClockInTime:   s.ClockInTime,
		ClockInLat:    s.ClockInLat,
		ClockInLng:    s.ClockInLng,
		ClockOutTime:  s.ClockOutTime,
		ClockOutLat:   s.ClockOutLat,
		ClockOutLng:   s.ClockOutLng,
		Location:      s.Location,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
