package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi"
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
// @Success		200		{object}	handlers.responseMsg
// @Failure		500		{object}	handlers.responseErr
// @Router		/people	[GET]
func HandleListPeople(logger *httplog.Logger, service peopleLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		people, err := service.ListPeople(ctx)
		if err != nil {
			logger.Error("error getting all people", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		peopleOut := mapMultipleOutput(people)
		encodeResponse(w, logger, http.StatusOK, responsePeople{
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
// @Param       id   path     int  true  "id"
// @Success		200		{object}	handlers.responseMsg
// @Failure		500		{object}	handlers.responseErr
// @Router      /people/{id}  [GET]
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
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}
		encodeResponse(w, logger, http.StatusOK, responsePerson{
			Person: mapOutput(person),
		})
	}
}