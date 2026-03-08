package errors

import "errors"

var (
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrTaskNotFound     = errors.New("task not found")
	ErrValidation       = errors.New("validation error")
	ErrInvalidStatus    = errors.New("invalid status transition")
)
