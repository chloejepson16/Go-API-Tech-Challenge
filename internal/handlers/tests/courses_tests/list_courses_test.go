package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers/courses"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCourseLister struct {
	mock.Mock
}

func (m *mockCourseLister) ListCourses(ctx context.Context) ([]models.Course, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *mockCourseLister) ListCourseByID(ctx context.Context, id int) (models.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Course), args.Error(1)
}

func TestHandleListCourses(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseLister)

	expectedCourses := []models.Course{
		{ID: 2, Name: "Databases"},
		{ID: 3, Name: "UI Design"},
		{ID: 1, Name: "Programmer"},
		{ID: 4, Name: "Introduction to Go"},
	}

	mockService.On("ListCourses", mock.Anything).Return(expectedCourses, nil)

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	rr := httptest.NewRecorder()
	handler := handlers.HandleListCourses(logger, mockService)
	handler.ServeHTTP(rr, req)

	expectedBody := `{
		"courses": [
			{"id": 2, "name": "Databases"},
			{"id": 3, "name": "UI Design"},
			{"id": 1, "name": "Programmer"},
			{"id": 4, "name": "Introduction to Go"}
		]
	}`

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, expectedBody, rr.Body.String())
}

// func TestHandleGetCourseById(t *testing.T) {
// 	logger := httplog.NewLogger("test", httplog.Options{})
// 	mockService := new(mockCourseLister)

// 	expectedCourse := models.Course{ID: 1, Name: "Programmer"}
// 	mockService.On("ListCourseByID", mock.Anything, 1).Return(expectedCourse, nil)
// 	req := httptest.NewRequest(http.MethodGet, "/courses/1", nil)
// 	rr := httptest.NewRecorder()

// 	handler := handlers.HandleGetCourseById(logger, mockService)
// 	handler.ServeHTTP(rr, req)
// 	expectedBody := `{
// 		"course": {
// 			"id": 1,
// 			"name": "Programmer"
// 		}
// 	}`

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	assert.JSONEq(t, expectedBody, rr.Body.String())
// }


