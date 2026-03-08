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
	"github.com/bsr144/evv-logger/impl/memory"
)

func setupScheduleApp() (*fiber.App, *memory.Store) {
	store := memory.New()
	uc := schedule.NewUsecase(store, store)
	ctrl := controller.NewScheduleController(uc)

	app := fiber.New()
	app.Post("/schedules", ctrl.Create)
	app.Get("/schedules", ctrl.GetAll)
	app.Get("/schedules/:id", ctrl.GetByID)
	app.Patch("/schedules/:id", ctrl.Update)

	return app, store
}

func TestScheduleController_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid schedule",
			body:           `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00","location":"Home"}`,
			expectedStatus: fiber.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "missing required fields",
			body:           `{"caregiver_name":"Alice"}`,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid json",
			body:           `{not valid json}`,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "empty body",
			body:           `{}`,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupScheduleApp()

			req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(tt.body))
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

func TestScheduleController_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		seedCount      int
		queryDate      string
		expectedStatus int
		expectedTotal  int
	}{
		{
			name:           "empty list",
			seedCount:      0,
			queryDate:      "",
			expectedStatus: fiber.StatusOK,
			expectedTotal:  0,
		},
		{
			name:           "multiple schedules",
			seedCount:      3,
			queryDate:      "",
			expectedStatus: fiber.StatusOK,
			expectedTotal:  3,
		},
		{
			name:           "filter by date",
			seedCount:      3,
			queryDate:      "2026-03-08",
			expectedStatus: fiber.StatusOK,
			expectedTotal:  3,
		},
		{
			name:           "filter by date no match",
			seedCount:      3,
			queryDate:      "2099-01-01",
			expectedStatus: fiber.StatusOK,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupScheduleApp()

			for i := 0; i < tt.seedCount; i++ {
				body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00"}`
				seedReq := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
				seedReq.Header.Set("Content-Type", "application/json")
				_, err := app.Test(seedReq, -1)
				require.NoError(t, err)
			}

			url := "/schedules"
			if tt.queryDate != "" {
				url += "?date=" + tt.queryDate
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var apiResp dto.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)
			assert.True(t, apiResp.Success)

			dataMap, ok := apiResp.Data.(map[string]interface{})
			require.True(t, ok)

			statsMap, ok := dataMap["stats"].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, float64(tt.expectedTotal), statsMap["total"])
		})
	}
}

func TestScheduleController_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		seed           bool
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "existing schedule",
			id:             "1",
			seed:           true,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name:           "not found",
			id:             "999",
			seed:           false,
			expectedStatus: fiber.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "invalid id",
			id:             "abc",
			seed:           false,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupScheduleApp()

			if tt.seed {
				body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00"}`
				seedReq := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
				seedReq.Header.Set("Content-Type", "application/json")
				_, err := app.Test(seedReq, -1)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodGet, "/schedules/"+tt.id, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var apiResp dto.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)

			if tt.expectedError {
				assert.False(t, apiResp.Success)
			} else {
				assert.True(t, apiResp.Success)
				assert.NotNil(t, apiResp.Data)
			}
		})
	}
}

func TestScheduleController_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           string
		seed           bool
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "valid update clock in",
			id:             "1",
			body:           `{"status":"in_progress","clock_in_lat":40.7128,"clock_in_lng":-74.0060}`,
			seed:           true,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name:           "not found",
			id:             "999",
			body:           `{"status":"in_progress"}`,
			seed:           false,
			expectedStatus: fiber.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "invalid id",
			id:             "abc",
			body:           `{"status":"in_progress"}`,
			seed:           false,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid status value",
			id:             "1",
			body:           `{"status":"nonexistent_status"}`,
			seed:           true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid json",
			id:             "1",
			body:           `{not valid}`,
			seed:           true,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _ := setupScheduleApp()

			if tt.seed {
				body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00"}`
				seedReq := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
				seedReq.Header.Set("Content-Type", "application/json")
				_, err := app.Test(seedReq, -1)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, "/schedules/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var apiResp dto.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			require.NoError(t, err)

			if tt.expectedError {
				assert.False(t, apiResp.Success)
			} else {
				assert.True(t, apiResp.Success)
			}
		})
	}
}

func TestScheduleController_Create_ResponseShape(t *testing.T) {
	app, _ := setupScheduleApp()

	body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00","location":"Home"}`
	req := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
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
	assert.Equal(t, "Alice", dataMap["caregiver_name"])
	assert.Equal(t, "Bob", dataMap["patient_name"])
	assert.Equal(t, "2026-03-08", dataMap["date"])
	assert.Equal(t, "upcoming", dataMap["status"])
}

func TestScheduleController_GetByID_ResponseShape(t *testing.T) {
	app, _ := setupScheduleApp()

	body := `{"caregiver_name":"Alice","patient_name":"Bob","date":"2026-03-08","start_time":"09:00","end_time":"17:00"}`
	seedReq := httptest.NewRequest(http.MethodPost, "/schedules", strings.NewReader(body))
	seedReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(seedReq, -1)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/schedules/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var apiResp dto.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	require.NoError(t, err)
	assert.True(t, apiResp.Success)

	dataMap, ok := apiResp.Data.(map[string]interface{})
	require.True(t, ok)

	schedMap, ok := dataMap["schedule"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(1), schedMap["id"])
	assert.Equal(t, "Alice", schedMap["caregiver_name"])

	tasksSlice, ok := dataMap["tasks"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, tasksSlice)
}
