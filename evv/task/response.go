package task

import "time"

type TaskResponse struct {
	ID          int       `json:"id"`
	ScheduleID  int       `json:"schedule_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
