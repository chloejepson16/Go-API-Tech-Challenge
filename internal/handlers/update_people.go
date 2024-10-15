package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/httplog/v2"
)

type personUpdater interface {
	UpdatePerson(ctx context.Context, id int, person models.Person) (models.Person, error)
}

// HandleUpdatePerson is a Handler that updates a user based on a user object from the request body.
//
// @Summary		Update a person by id
// @Description	Update a person by id
// @Tags		people
// @Accept		json
// @Produce		json
// @Param		id			path		int	true						"id"
// @Param		person		body		handlers.InputPerson		true	"Person Object"
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/people/{id}	[PUT]
func HandleUpdatePerson(logger *httplog.Logger, service personUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get and validate ID
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			EncodeResponse(w, logger, http.StatusBadRequest, ResponseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get and validate body as object
		personIn, problems, err := DecodeValidateBody[InputPerson, models.Person](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				EncodeResponse(w, logger, http.StatusBadRequest, ResponseErr{
					Error: "error validating personIn",
				})
			default:
				logger.Error("BodyParser error", "error", err)
				EncodeResponse(w, logger, http.StatusBadRequest, ResponseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		// update object in database
		user, err := service.UpdatePerson(ctx, id, personIn)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			EncodeResponse(w, logger, http.StatusInternalServerError, ResponseErr{
				Error: "Error updating object",
			})
			return
		}

		// return response
		personOut := MapOutput(user)
		EncodeResponse(w, logger, http.StatusOK, ResponsePerson{
			Person: personOut,
		})
	}
}