package handlers_test

import (
	"context"
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

type mockCourseDeleter struct {
	mock.Mock
}

func (m *mockCourseDeleter) DeleteCourseByID(ctx context.Context, id int) (models.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Course), args.Error(1)
}

func TestHandleDeleteCourseByID(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockCourseDeleter)

	expectedCourse := models.Course{ID: 1, Name: "Programmer"}
	mockService.On("DeleteCourseByID", mock.Anything, expectedCourse.ID).Return(expectedCourse, nil)

	req := httptest.NewRequest(http.MethodDelete, "/courses/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(expectedCourse.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler := handlers.HandleDeleteCourseByID(logger, mockService)
	handler.ServeHTTP(rr, req)

	expectedBody := `{
		"course": {
			"id": 1,
			"name": "Programmer"
		}
	}`

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, expectedBody, rr.Body.String())
	mockService.AssertCalled(t, "DeleteCourseByID", mock.Anything, expectedCourse.ID)
}
