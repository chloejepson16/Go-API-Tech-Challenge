package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type peopleLister interface{
	ListPeople(ctx context.Context) ([]models.Person, error)
	ListPersonByID(ctx context.Context, id int) (models.Person, error)
}

// HandleListPeople is a Handler that returns a list of all people.
//
// @Summary		List all people
// @Description	List all people
// @Tags		people
// @Accept		json
// @Produce		json
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/people	[GET]
func HandleListPeople(logger *httplog.Logger, service peopleLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		people, err := service.ListPeople(ctx)
		if err != nil {
			logger.Error("error getting all people", "error", err)
			handlers.EncodeResponse(w, logger, http.StatusInternalServerError, handlers.ResponseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		peopleOut := handlers.MapMultipleOutput(people)
		handlers.EncodeResponse(w, logger, http.StatusOK, handlers.ResponsePeople{
			People: peopleOut,
		})
	}
}

// HandleGetPersonByID is a Handler that returns a single person by ID.
//
// @Summary     Get person by ID
// @Description Get a person by their ID
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       id      path    int  true  "id"
// @Success     200     {object} handlers.ResponseMsg
// @Failure     500     {object} handlers.ResponseErr
// @Router      /people/{id} [GET]
func HandleGetPersonByID(logger *httplog.Logger, service peopleLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()
		personId := chi.URLParam(r, "id")
		fmt.Printf("ID is: %s", personId)
		id, err := strconv.Atoi(personId)
		fmt.Printf("ID 2 is: %d", id)

		// get values from database
		person, err := service.ListPersonByID(ctx, id)
		if err != nil {
			logger.Error("error getting specific people", "error", err)
			handlers.EncodeResponse(w, logger, http.StatusInternalServerError, handlers.ResponseErr{
				Error: "Error retrieving data",
			})
			return
		}
		handlers.EncodeResponse(w, logger, http.StatusOK, handlers.ResponsePerson{
			Person: handlers.MapOutput(person),
		})
	}
}