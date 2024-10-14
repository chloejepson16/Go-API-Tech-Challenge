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

type peopleDeleter interface{
	DeletePersonByID(ctx context.Context, id int) (models.Person, error)
}

// HandleDeletePersonByID is a Handler that deletes a single person by ID.
//
// @Summary     Delete person by ID
// @Description Delete a person by their ID
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "id"
// @Success     200  {object}    handlers.ResponseMsg
// @Failure     500  {object}    handlers.ResponseErr
// @Router      /people/{id}  [DELETE]
func HandleDeletePersonByID(logger *httplog.Logger, service peopleDeleter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		// setup
		ctx := r.Context()
		personId := chi.URLParam(r, "id")
		fmt.Printf("ID is: %s", personId)
		id, err := strconv.Atoi(personId)
		fmt.Printf("ID 2 is: %d", id)

		person, err:= service.DeletePersonByID(ctx, id)
		if err != nil{
			logger.Error("error deleting specific people", "error", err)
			EncodeResponse(w, logger, http.StatusInternalServerError, ResponseErr{
				Error: "Error retrieving data",
			})
			return
		}
		EncodeResponse(w, logger, http.StatusOK, ResponsePerson{
			Person: MapOutput(person),
		})
	}
}