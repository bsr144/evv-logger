package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/bsr144/evv-logger/delivery/http/controller"
	"github.com/bsr144/evv-logger/delivery/http/middleware"
	"github.com/bsr144/evv-logger/delivery/http/router"
	"github.com/bsr144/evv-logger/entity"
	"github.com/bsr144/evv-logger/evv/schedule"
	"github.com/bsr144/evv-logger/evv/task"
	"github.com/bsr144/evv-logger/impl/memory"
)

func main() {
	store := memory.New()
	seedData(store)

	scheduleUC := schedule.NewUsecase(store, store)
	taskUC := task.NewUsecase(store, store)

	scheduleCtrl := controller.NewScheduleController(scheduleUC)
	taskCtrl := controller.NewTaskController(taskUC)

	app := fiber.New(fiber.Config{
		AppName:      "EVV Logger",
		BodyLimit:    64 * 1024,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})
	app.Use(recover.New())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	router.Setup(app, scheduleCtrl, taskCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting EVV Logger server", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func seedData(store *memory.Store) {
	ctx := context.Background()
	now := time.Now()

	schedules := []*entity.Schedule{
		{
			CaregiverName: "Alice Johnson",
			PatientName:   "Robert Smith",
			Date:          "2026-03-08",
			StartTime:     "08:00",
			EndTime:       "12:00",
			Status:        entity.ScheduleStatusUpcoming,
			Location:      "123 Main St",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			CaregiverName: "Alice Johnson",
			PatientName:   "Mary Williams",
			Date:          "2026-03-08",
			StartTime:     "13:00",
			EndTime:       "17:00",
			Status:        entity.ScheduleStatusUpcoming,
			Location:      "456 Oak Ave",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			CaregiverName: "Bob Davis",
			PatientName:   "James Brown",
			Date:          "2026-03-08",
			StartTime:     "09:00",
			EndTime:       "15:00",
			Status:        entity.ScheduleStatusInProgress,
			Location:      "789 Pine Rd",
			ClockInTime:   &now,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			CaregiverName: "Carol White",
			PatientName:   "Patricia Jones",
			Date:          "2026-03-07",
			StartTime:     "10:00",
			EndTime:       "14:00",
			Status:        entity.ScheduleStatusCompleted,
			Location:      "321 Elm Blvd",
			ClockInTime:   &now,
			ClockOutTime:  &now,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			CaregiverName: "David Garcia",
			PatientName:   "Linda Martinez",
			Date:          "2026-03-07",
			StartTime:     "07:00",
			EndTime:       "11:00",
			Status:        entity.ScheduleStatusMissed,
			Location:      "654 Birch Ln",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	scheduleIDs := make([]int, 0, len(schedules))
	for _, s := range schedules {
		created, err := store.CreateSchedule(ctx, s)
		if err != nil {
			slog.Error("failed to seed schedule", "error", err)
			continue
		}
		scheduleIDs = append(scheduleIDs, created.ID)
	}

	tasks := []*entity.Task{
		{ScheduleID: scheduleIDs[0], Description: "Check vital signs", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[0], Description: "Administer medication", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[1], Description: "Physical therapy exercises", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[1], Description: "Wound care", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[2], Description: "Blood pressure check", Status: entity.TaskStatusCompleted, Notes: "BP 120/80", CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[2], Description: "Meal preparation", Status: entity.TaskStatusPending, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[2], Description: "Mobility assistance", Status: entity.TaskStatusNotCompleted, Notes: "Patient too fatigued", CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[3], Description: "Medication review", Status: entity.TaskStatusCompleted, CreatedAt: now, UpdatedAt: now},
		{ScheduleID: scheduleIDs[3], Description: "Health assessment", Status: entity.TaskStatusCompleted, Notes: "All clear", CreatedAt: now, UpdatedAt: now},
	}

	for _, t := range tasks {
		if _, err := store.CreateTask(ctx, t); err != nil {
			slog.Error("failed to seed task", "error", err)
		}
	}

	slog.Info("seed data loaded", "schedules", len(scheduleIDs), "tasks", len(tasks))
}
