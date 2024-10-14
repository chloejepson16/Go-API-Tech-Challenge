package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCourseCreator struct {
	mock.Mock
}

func (m *mockCourseCreator) CreateCourse(ctx context.Context, course models.Course) (models.Course, error) {
	args := m.Called(ctx, course)
	return args.Get(0).(models.Course), args.Error(1)
}

func TestHandleCreateCourse(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseCreator)

	courseIn := models.Course{ID: 0, Name: "New Course"}
	expectedCourse := models.Course{ID: 0, Name: "New Course"}

	mockService.On("CreateCourse", mock.Anything, courseIn).Return(expectedCourse, nil)
	body, _ := json.Marshal(courseIn)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := handlers.HandleCreateCourse(logger, mockService)
	handler.ServeHTTP(rr, req)

	expectedBody := `{
		"course": {
			"id": 0,
			"name": "New Course"
		}
	}`

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, expectedBody, rr.Body.String())
}
