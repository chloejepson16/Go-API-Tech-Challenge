package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
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
			handlers.EncodeResponse(w, logger, http.StatusBadRequest, handlers.ResponseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get and validate body as object
		personIn, problems, err := handlers.DecodeValidateBody[handlers.InputPerson, models.Person](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				handlers.EncodeResponse(w, logger, http.StatusBadRequest, handlers.ResponseErr{
					Error: "error validating personIn",
				})
			default:
				logger.Error("BodyParser error", "error", err)
				handlers.EncodeResponse(w, logger, http.StatusBadRequest, handlers.ResponseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		// update object in database
		user, err := service.UpdatePerson(ctx, id, personIn)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			handlers.EncodeResponse(w, logger, http.StatusInternalServerError, handlers.ResponseErr{
				Error: "Error updating object",
			})
			return
		}

		// return response
		personOut := handlers.MapOutput(user)
		handlers.EncodeResponse(w, logger, http.StatusOK, handlers.ResponsePerson{
			Person: personOut,
		})
	}
}