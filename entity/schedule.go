package entity

import "time"

type Schedule struct {
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

const (
	ScheduleStatusUpcoming   = "upcoming"
	ScheduleStatusInProgress = "in_progress"
	ScheduleStatusCompleted  = "completed"
	ScheduleStatusMissed     = "missed"
)
