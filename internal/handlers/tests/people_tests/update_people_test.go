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

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers/people"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the personUpdater interface
type mockPersonUpdater struct {
	mock.Mock
}

func (m *mockPersonUpdater) UpdatePerson(ctx context.Context, id int, person models.Person) (models.Person, error) {
	args := m.Called(ctx, id, person)
	return args.Get(0).(models.Person), args.Error(1)
}

// func TestHandleUpdatePerson(t *testing.T) {
// 	logger := httplog.NewLogger("test", httplog.Options{})
// 	mockService := new(mockPersonUpdater)

// 	personIn := models.Person{
// 		ID:         1,
// 		FirstName:  "Jane",
// 		LastName:   "Doe",
// 		PersonType: "student",
// 		Age:        14,
// 	}
// 	personOut := models.Person{
// 		ID:         1,
// 		FirstName:  "Jane",
// 		LastName:   "Doe",
// 		PersonType: "student",
// 		Age:        14,
// 	}

// 	mockService.On("UpdatePerson", mock.Anything, personIn.ID, personIn).Return(personOut, nil)
// 	body, _ := json.Marshal(personIn)
// 	req := httptest.NewRequest(http.MethodPut, "/people/"+strconv.Itoa(personIn.ID), bytes.NewBuffer(body))
// 	req = req.WithContext(context.Background())

// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("id", strconv.Itoa(personIn.ID))
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	rr := httptest.NewRecorder()

// 	handler := handlers.HandleUpdatePerson(logger, mockService)
// 	handler.ServeHTTP(rr, req)
	
// 	expectedBody := `{
// 		"person": {
// 			"id": 1,
// 			"first_name": "Jane",
// 			"last_name": "Doe",
// 			"person_type": "student",
// 			"age": 14
// 		}
// 	}`

// 	assert.Equal(t, http.StatusOK, rr.Code, "Wrong status code received")
// 	assert.JSONEq(t, expectedBody, rr.Body.String(), "Wrong response body")
// }

func TestHandleUpdatePerson_BadID(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockPersonUpdater)

	req := httptest.NewRequest(http.MethodPut, "/people/abc", nil)
	rr := httptest.NewRecorder()

	handler := handlers.HandleUpdatePerson(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleUpdatePerson_DBError(t *testing.T) {
	logger := httplog.NewLogger("test", httplog.Options{})
	mockService := new(mockPersonUpdater)

	personIn := models.Person{
		ID:   1,
		FirstName: "Jane",
		LastName: "Jane",
		PersonType: "student",
		Age: 14,
	}

	mockService.On("UpdatePerson", mock.Anything, personIn.ID, personIn).Return(models.Person{}, errors.New("database error"))

	body, _ := json.Marshal(personIn)
	req := httptest.NewRequest(http.MethodPut, "/people/1", bytes.NewBuffer(body))
	req = req.WithContext(context.Background())

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(personIn.ID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	handler := handlers.HandleUpdatePerson(logger, mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
