package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bsr144/evv-logger/delivery/http/controller"
	"github.com/bsr144/evv-logger/delivery/http/dto"
	"github.com/bsr144/evv-logger/evv/schedule"
	"github.com/bsr144/evv-logger/evv/task"
	"github.com/bsr144/evv-logger/impl/memory"
)

func setupTaskApp() (*fiber.App, *memory.Store) {
	store := memory.New()
	schedUC := schedule.NewUsecase(store, store)
	taskUC := task.NewUsecase(store, store)
	schedCtrl := controller.NewScheduleController(schedUC)
	taskCtrl := controller.NewTaskController(taskUC)

	app := fiber.New()
	app.Post("/schedules", schedCtrl.Create)
	app.Post("/schedules/:id/tasks", taskCtrl.Create)
	app.Get("/schedules/:id/tasks", taskCtrl.GetByScheduleID)
	app.Patch("/schedules/:id/tasks/:taskId", taskCtrl.Update)

	return app, store
}

func seedSchedule(t *testing.T, app *fiber.App) {
	t.Helper()
	body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00"}`
	req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func seedTask(t *testing.T, app *fiber.App, scheduleID string) {
	t.Helper()
	body := `{"description":"Check vitals"}`
	req := httptest.NewRequest(http.MethodPost, "/schedules/"+scheduleID+"/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestTaskController_Create(t *testing.T) {
	tests := []struct {
		name           string
		scheduleID     string
		body           string
		seedSchedule   bool
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid task",
			scheduleID:     "1",
			body:           `{"description":"Check vitals"}`,
			seedSchedule:   true,
			expectedStatus: fiber.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "schedule not found",
			scheduleID:     "999",
			body:           `{"description":"Check vitals"}`,
			seedSchedule:   false,
			expectedStatus: fiber.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "missing description",
			scheduleID:     "1",
			body:           `{}`,
			seedSchedule:   true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid json",
			scheduleID:     "1",
			body:           `{invalid}`,
			seedSchedule:   true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid schedule id",
			scheduleID:     "abc",
			body:           `{"description":"Check vitals"}`,
			seedSchedule:   false,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupTaskApp()

			if tt.seedSchedule {
				seedSchedule(t, app)
			}

			req := httptest.NewRequest(http.MethodPost, "/schedules/"+tt.scheduleID+"/tasks", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var apiResp dto.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)

			if tt.expectedError {
				assert.False(t, apiResp.Success)
				assert.NotEmpty(t, apiResp.Error)
			} else {
				assert.True(t, apiResp.Success)
				assert.NotNil(t, apiResp.Data)
			}
		})
	}
}

func TestTaskController_GetByScheduleID(t *testing.T) {
	tests := []struct {
		name           string
		scheduleID     string
		seedSchedule   bool
		seedTaskCount  int
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "tasks found",
			scheduleID:     "1",
			seedSchedule:   true,
			seedTaskCount:  2,
			expectedStatus: fiber.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "no tasks",
			scheduleID:     "1",
			seedSchedule:   true,
			seedTaskCount:  0,
			expectedStatus: fiber.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "schedule not found",
			scheduleID:     "999",
			seedSchedule:   false,
			seedTaskCount:  0,
			expectedStatus: fiber.StatusNotFound,
			expectedCount:  0,
		},
		{
			name:           "invalid schedule id",
			scheduleID:     "abc",
			seedSchedule:   false,
			seedTaskCount:  0,
			expectedStatus: fiber.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupTaskApp()

			if tt.seedSchedule {
				seedSchedule(t, app)
			}
			for i := 0; i < tt.seedTaskCount; i++ {
				seedTask(t, app, tt.scheduleID)
			}

			req := httptest.NewRequest(http.MethodGet, "/schedules/"+tt.scheduleID+"/tasks", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == fiber.StatusOK {
				var apiResp dto.APIResponse
				err = json.NewDecoder(resp.Body).Decode(&apiResp)
				require.NoError(t, err)
				assert.True(t, apiResp.Success)

				dataSlice, ok := apiResp.Data.([]interface{})
				require.True(t, ok)
				assert.Len(t, dataSlice, tt.expectedCount)
			}
		})
	}
}

func TestTaskController_Update(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		body           string
		seedSchedule   bool
		seedTask       bool
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "mark completed",
			taskID:         "1",
			body:           `{"status":"completed"}`,
			seedSchedule:   true,
			seedTask:       true,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name:           "not_completed with notes",
			taskID:         "1",
			body:           `{"status":"not_completed","notes":"Patient refused"}`,
			seedSchedule:   true,
			seedTask:       true,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name:           "not_completed without notes",
			taskID:         "1",
			body:           `{"status":"not_completed"}`,
			seedSchedule:   true,
			seedTask:       true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "task not found",
			taskID:         "999",
			body:           `{"status":"completed"}`,
			seedSchedule:   false,
			seedTask:       false,
			expectedStatus: fiber.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "invalid task id",
			taskID:         "abc",
			body:           `{"status":"completed"}`,
			seedSchedule:   false,
			seedTask:       false,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid json",
			taskID:         "1",
			body:           `{invalid}`,
			seedSchedule:   true,
			seedTask:       true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "update notes only",
			taskID:         "1",
			body:           `{"notes":"Some observations"}`,
			seedSchedule:   true,
			seedTask:       true,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupTaskApp()

			if tt.seedSchedule {
				seedSchedule(t, app)
			}
			if tt.seedTask {
				seedTask(t, app, "1")
			}

			req := httptest.NewRequest(http.MethodPatch, "/schedules/1/tasks/"+tt.taskID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var apiResp dto.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)

			if tt.expectedError {
				assert.False(t, apiResp.Success)
				assert.NotEmpty(t, apiResp.Error)
			} else {
				assert.True(t, apiResp.Success)
			}
		})
	}
}

func TestTaskController_Create_ResponseShape(t *testing.T) {
	app, _ := setupTaskApp()
	seedSchedule(t, app)

	body := `{"description":"Administer medication"}`
	req := httptest.NewRequest(http.MethodPost, "/schedules/1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var apiResp dto.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	require.NoError(t, err)
	assert.True(t, apiResp.Success)

	dataMap, ok := apiResp.Data.(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, float64(1), dataMap["id"])
	assert.Equal(t, float64(1), dataMap["schedule_id"])
	assert.Equal(t, "Administer medication", dataMap["description"])
	assert.Equal(t, "pending", dataMap["status"])
}
