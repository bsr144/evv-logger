package task

type CreateTaskRequest struct {
	Description string `json:"description" validate:"required,max=500"`
}

type UpdateTaskRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=pending completed not_completed"`
	Notes  *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}
