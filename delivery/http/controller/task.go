package controller

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/bsr144/evv-logger/delivery/http/dto"
	"github.com/bsr144/evv-logger/evv/task"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

type TaskController struct {
	usecase  task.Usecase
	validate *validator.Validate
}

func NewTaskController(usecase task.Usecase) *TaskController {
	return &TaskController{usecase: usecase, validate: validator.New()}
}

func (tc *TaskController) Create(c *fiber.Ctx) error {
	scheduleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid schedule id"))
	}
	if scheduleID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("schedule id must be positive"))
	}

	var req task.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid request body"))
	}

	if err := tc.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse(err))
	}

	resp, err := tc.usecase.Create(c.Context(), scheduleID, &req)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrScheduleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("schedule not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to create task"))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse(resp))
}

func (tc *TaskController) GetByScheduleID(c *fiber.Ctx) error {
	scheduleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid schedule id"))
	}
	if scheduleID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("schedule id must be positive"))
	}

	resp, err := tc.usecase.GetByScheduleID(c.Context(), scheduleID)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrScheduleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("schedule not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to get tasks"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}

func (tc *TaskController) Update(c *fiber.Ctx) error {
	scheduleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid schedule id"))
	}
	if scheduleID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("schedule id must be positive"))
	}

	taskID, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid task id"))
	}
	if taskID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("task id must be positive"))
	}

	var req task.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid request body"))
	}

	resp, err := tc.usecase.Update(c.Context(), scheduleID, taskID, &req)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrTaskNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("task not found"))
		}
		if errors.Is(err, pkgerrors.ErrScheduleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("schedule not found"))
		}
		if errors.Is(err, pkgerrors.ErrValidation) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("notes required when status is not_completed"))
		}
		if errors.Is(err, pkgerrors.ErrInvalidStatus) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid status"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to update task"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}
