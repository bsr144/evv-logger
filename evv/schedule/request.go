package schedule

import "time"

type CreateScheduleRequest struct {
	CaregiverName string `json:"caregiver_name" validate:"required,max=200"`
	PatientName   string `json:"patient_name" validate:"required,max=200"`
	Date          string `json:"date" validate:"required,max=10,datetime=2006-01-02"`
	StartTime     string `json:"start_time" validate:"required,max=5"`
	EndTime       string `json:"end_time" validate:"required,max=5"`
	Location      string `json:"location" validate:"max=500"`
}

type UpdateScheduleRequest struct {
	Status       *string    `json:"status,omitempty"`
	ClockInTime  *time.Time `json:"clock_in_time,omitempty"`
	ClockInLat   *float64   `json:"clock_in_lat,omitempty"`
	ClockInLng   *float64   `json:"clock_in_lng,omitempty"`
	ClockOutTime *time.Time `json:"clock_out_time,omitempty"`
	ClockOutLat  *float64   `json:"clock_out_lat,omitempty"`
	ClockOutLng  *float64   `json:"clock_out_lng,omitempty"`
}
