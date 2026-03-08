package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/bsr144/evv-logger/entity"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
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

func (s *Store) CreateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	schedule.ID = s.scheduleNextID
	s.scheduleNextID++

	stored := deepCopySchedule(schedule)
	s.schedules[stored.ID] = stored

	return deepCopySchedule(stored), nil
}

func (s *Store) GetScheduleByID(_ context.Context, id int) (*entity.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedule, ok := s.schedules[id]
	if !ok {
		return nil, pkgerrors.ErrScheduleNotFound
	}

	return deepCopySchedule(schedule), nil
}

func (s *Store) GetAllSchedules(_ context.Context) ([]*entity.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*entity.Schedule, 0, len(s.schedules))
	for _, schedule := range s.schedules {
		result = append(result, deepCopySchedule(schedule))
	}

	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })

	return result, nil
}

func (s *Store) UpdateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.schedules[schedule.ID]; !ok {
		return nil, pkgerrors.ErrScheduleNotFound
	}

	stored := deepCopySchedule(schedule)
	s.schedules[stored.ID] = stored

	return deepCopySchedule(stored), nil
}

func (s *Store) CreateTask(_ context.Context, task *entity.Task) (*entity.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.taskNextID
	s.taskNextID++

	stored := *task
	s.tasks[stored.ID] = &stored

	result := stored
	return &result, nil
}

func (s *Store) GetTaskByID(_ context.Context, id int) (*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, pkgerrors.ErrTaskNotFound
	}

	copy := *task
	return &copy, nil
}

func (s *Store) GetTasksByScheduleID(_ context.Context, scheduleID int) ([]*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*entity.Task, 0)
	for _, task := range s.tasks {
		if task.ScheduleID == scheduleID {
			cp := *task
			result = append(result, &cp)
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })

	return result, nil
}

func (s *Store) UpdateTask(_ context.Context, task *entity.Task) (*entity.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return nil, pkgerrors.ErrTaskNotFound
	}

	stored := *task
	s.tasks[stored.ID] = &stored

	result := stored
	return &result, nil
}

func deepCopySchedule(src *entity.Schedule) *entity.Schedule {
	cp := *src
	if src.ClockInTime != nil {
		t := *src.ClockInTime
		cp.ClockInTime = &t
	}
	if src.ClockInLat != nil {
		v := *src.ClockInLat
		cp.ClockInLat = &v
	}
	if src.ClockInLng != nil {
		v := *src.ClockInLng
		cp.ClockInLng = &v
	}
	if src.ClockOutTime != nil {
		t := *src.ClockOutTime
		cp.ClockOutTime = &t
	}
	if src.ClockOutLat != nil {
		v := *src.ClockOutLat
		cp.ClockOutLat = &v
	}
	if src.ClockOutLng != nil {
		v := *src.ClockOutLng
		cp.ClockOutLng = &v
	}
	return &cp
}

var _ interface {
	CreateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	GetScheduleByID(_ context.Context, id int) (*entity.Schedule, error)
	GetAllSchedules(_ context.Context) ([]*entity.Schedule, error)
	UpdateSchedule(_ context.Context, schedule *entity.Schedule) (*entity.Schedule, error)
	CreateTask(_ context.Context, task *entity.Task) (*entity.Task, error)
	GetTaskByID(_ context.Context, id int) (*entity.Task, error)
	GetTasksByScheduleID(_ context.Context, scheduleID int) ([]*entity.Task, error)
	UpdateTask(_ context.Context, task *entity.Task) (*entity.Task, error)
} = (*Store)(nil)
