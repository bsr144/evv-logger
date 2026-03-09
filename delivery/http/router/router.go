package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bsr144/evv-logger/delivery/http/controller"
)

func Setup(app *fiber.App, sc *controller.ScheduleController, tc *controller.TaskController) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api.Post("/schedules", sc.Create)
	api.Get("/schedules", sc.GetAll)
	api.Get("/schedules/stats", sc.Stats)
	api.Get("/schedules/:id", sc.GetByID)
	api.Patch("/schedules/:id", sc.Update)

	api.Post("/schedules/:id/tasks", tc.Create)
	api.Get("/schedules/:id/tasks", tc.GetByScheduleID)
	api.Patch("/schedules/:id/tasks/:taskId", tc.Update)
}
