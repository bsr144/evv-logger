package controller

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/bsr144/evv-logger/delivery/http/dto"
	"github.com/bsr144/evv-logger/evv/schedule"
	pkgerrors "github.com/bsr144/evv-logger/pkg/errors"
)

type ScheduleController struct {
	usecase  schedule.Usecase
	validate *validator.Validate
}

func NewScheduleController(usecase schedule.Usecase) *ScheduleController {
	return &ScheduleController{usecase: usecase, validate: validator.New()}
}

func (sc *ScheduleController) Create(c *fiber.Ctx) error {
	var req schedule.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid request body"))
	}

	if err := sc.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse(err))
	}

	resp, err := sc.usecase.Create(c.Context(), &req)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrValidation) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse(err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to create schedule"))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse(resp))
}

func (sc *ScheduleController) GetAll(c *fiber.Ctx) error {
	dateFilter := c.Query("date")

	resp, err := sc.usecase.GetAll(c.Context(), dateFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to get schedules"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}

func (sc *ScheduleController) Stats(c *fiber.Ctx) error {
	resp, err := sc.usecase.GetStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to get schedule stats"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}

func (sc *ScheduleController) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid schedule id"))
	}
	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("schedule id must be positive"))
	}

	resp, err := sc.usecase.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrScheduleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("schedule not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to get schedule"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}

func (sc *ScheduleController) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid schedule id"))
	}
	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("schedule id must be positive"))
	}

	var req schedule.UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid request body"))
	}

	resp, err := sc.usecase.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, pkgerrors.ErrScheduleNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse("schedule not found"))
		}
		if errors.Is(err, pkgerrors.ErrInvalidStatus) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse("invalid status transition"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse("failed to update schedule"))
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse(resp))
}
