package schedule

import "time"

type ScheduleResponse struct {
	ID            int        `json:"id"`
	CaregiverName string     `json:"caregiver_name"`
	PatientName   string     `json:"patient_name"`
	Date          string     `json:"date"`
	StartTime     string     `json:"start_time"`
	EndTime       string     `json:"end_time"`
	Status        string     `json:"status"`
	ClockInTime   *time.Time `json:"clock_in_time"`
	ClockInLat    *float64   `json:"clock_in_lat"`
	ClockInLng    *float64   `json:"clock_in_lng"`
	ClockOutTime  *time.Time `json:"clock_out_time"`
	ClockOutLat   *float64   `json:"clock_out_lat"`
	ClockOutLng   *float64   `json:"clock_out_lng"`
	Location      string     `json:"location"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type TaskSummary struct {
	ID          int       `json:"id"`
	ScheduleID  int       `json:"schedule_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ScheduleStats struct {
	Total      int `json:"total"`
	Upcoming   int `json:"upcoming"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Missed     int `json:"missed"`
}

type ScheduleListResponse struct {
	Data  []*ScheduleResponse `json:"data"`
	Stats ScheduleStats       `json:"stats"`
}

type ScheduleDetailResponse struct {
	Schedule *ScheduleResponse `json:"schedule"`
	Tasks    []*TaskSummary    `json:"tasks"`
}
