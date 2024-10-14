package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the courseUpdater interface
type mockCourseUpdater struct {
	mock.Mock
}

func (m *mockCourseUpdater) UpdateCourse(ctx context.Context, id int, course models.Course) (models.Course, error) {
	args := m.Called(ctx, id, course)
	return args.Get(0).(models.Course), args.Error(1)
}

func TestHandleUpdateCourse(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseUpdater)

	courseIn := models.Course{
		ID:    12,
		Name: "course 12",
	}
	courseOut := models.Course{
		ID:    12,
		Name: "course 12",
	}

	mockService.On("UpdateCourse", mock.Anything, courseIn.ID, courseIn).Return(courseOut, nil)
	body, _ := json.Marshal(courseIn)
	req := httptest.NewRequest(http.MethodPut, "/courses/12", bytes.NewBuffer(body))
	req = req.WithContext(context.Background())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(courseIn.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	handler := handlers.HandleUpdateCourse(logger, mockService)
	handler.ServeHTTP(rr, req)

	expectedBody := `{
		"course": {
			"id": 12,
			"name": "course 12"
		}
	}`

	t.Log("This is a test log:", rr.Body.String())


	assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code received")
	assert.JSONEq(t, expectedBody, rr.Body.String(), "Wrong response body")
}

func TestHandleUpdateCourse_BadID(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseUpdater)

	req := httptest.NewRequest(http.MethodPut, "/courses/abc", nil)
	rr := httptest.NewRecorder()

	handler := handlers.HandleUpdateCourse(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleUpdateCourse_DBError(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseUpdater)

	courseIn := models.Course{
		ID:    1,
		Name: "Test Course",
	}

	mockService.On("UpdateCourse", mock.Anything, courseIn.ID, courseIn).Return(models.Course{}, errors.New("database error"))

	body, _ := json.Marshal(courseIn)
	req := httptest.NewRequest(http.MethodPut, "/courses/1", bytes.NewBuffer(body))
	req = req.WithContext(context.Background())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(courseIn.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	handler := handlers.HandleUpdateCourse(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
