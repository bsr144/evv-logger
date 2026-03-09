package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{Success: true, Data: data}
}

func ErrorResponse(msg string) APIResponse {
	return APIResponse{Success: false, Error: msg}
}

func ValidationErrorResponse(err error) APIResponse {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse("validation failed")
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		msgs = append(msgs, fmt.Sprintf("'%s' failed on '%s'", strings.ToLower(fe.Field()), fe.Tag()))
	}

	return APIResponse{Success: false, Error: "validation failed: " + strings.Join(msgs, ", ")}
}
