package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the peopleLister interface
type mockPeopleLister struct {
	mock.Mock
}

func (m *mockPeopleLister) ListPeople(ctx context.Context) ([]models.Person, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (m *mockPeopleLister) ListPersonByID(ctx context.Context, id int) (models.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Person), args.Error(1)
}

func TestHandleListPeople(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockPeopleLister)

	people := []models.Person{
		{ID: 2, FirstName: "Jeff", LastName: "Bezos", PersonType: "professor", Age: 60},
		{ID: 3, FirstName: "Larry", LastName: "Page", PersonType: "student", Age: 51},
		{ID: 4, FirstName: "Bill", LastName: "Gates", PersonType: "student", Age: 67},
		{ID: 5, FirstName: "Elon", LastName: "Musk", PersonType: "student", Age: 52},
		{ID: 1, FirstName: "Steves", LastName: "JOBS", PersonType: "professor", Age: 12},
		{ID: 99, FirstName: "New", LastName: "Person", PersonType: "student", Age: 12},
		{ID: 13, FirstName: "John", LastName: "Doe", PersonType: "student", Age: 25},
		{ID: 14, FirstName: "John", LastName: "Doe", PersonType: "student", Age: 25},
	}

	req := httptest.NewRequest(http.MethodGet, "/people", nil)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()

	handler := handlers.HandleListPeople(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := handlers.ResponsePeople{
		People: handlers.MapMultipleOutput(people),
	}
	expectedBody, err := json.Marshal(expectedResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedBody), rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandleGetPersonByID(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockPeopleLister)

	expectedPerson := models.Person{
		ID:        13,
		FirstName: "John",
		LastName:  "Doe",
		PersonType: "student",
		Age:       25,
	}

	mockService.On("ListPersonByID", mock.Anything, expectedPerson.ID).Return(expectedPerson, nil)
	req := httptest.NewRequest(http.MethodGet, "/people/"+strconv.Itoa(expectedPerson.ID), nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(expectedPerson.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := handlers.HandleGetPersonByID(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response handlers.ResponsePerson
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedBody := handlers.ResponsePerson{
		Person: handlers.OutputPerson{
			ID:        expectedPerson.ID,
			FirstName: expectedPerson.FirstName,
			LastName:  expectedPerson.LastName,
			PersonType: expectedPerson.PersonType,
			Age:       expectedPerson.Age,
		},
	}

	assert.Equal(t, expectedBody, response)
	mockService.AssertExpectations(t)
}