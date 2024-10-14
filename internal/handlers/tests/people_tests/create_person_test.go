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
	"github.com/stretchr/testify/mock"
)

// Mock for the personCreator interface
type mockPersonCreator struct {
	mock.Mock
}

type OutputPerson struct {
	ID        int   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PersonType string `json:"type"`
	Age       int    `json:"age"`
}

type ResponsePerson struct{
	Person OutputPerson `json:"person"`
}

func (m *mockPersonCreator) CreatePerson(ctx context.Context, person models.Person) (models.Person, error) {
	args := m.Called(ctx, person)
	return args.Get(0).(models.Person), args.Error(1)
}

func TestHandleCreatePerson(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockPersonCreator)

	personIn := models.Person{
		ID:        15,
		FirstName: "John",
		LastName:  "Doe",
		PersonType: "student",
		Age:       25,
	}

	expectedPerson := models.Person{
		ID:        15,
		FirstName: "John",
		LastName:  "Doe",
		PersonType: "student",
		Age:       25,
	}

	mockService.On("CreatePerson", mock.Anything, personIn).Return(expectedPerson, nil)

	body, _ := json.Marshal(personIn)
	req := httptest.NewRequest(http.MethodPost, "/people", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	handler := handlers.HandleCreatePerson(logger, mockService)
	handler.ServeHTTP(rr, req)
	// assert.Equal(t, http.StatusOK, rr.Code)

	// expectedBody := `{
	// 	"person": {
	// 		"id": 15,
	// 		"first_name": "John",
	// 		"last_name": "Doe",
	// 		"type": "student",
	// 		"age": 25
	// 	}
	// }`

	
	// assert.JSONEq(t, expectedBody, rr.Body.String())
	// mockService.AssertExpectations(t)
}