package handlers

import (
	"context"
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
)

type personCreator interface{
	CreatePerson(ctx context.Context, person models.Person) (models.Person, error)
}

// HandleCreatePerson is a Handler that updates a user based on a user object from the request body.
//
// @Summary		Create a person
// @Description	Create a person
// @Tags		people
// @Accept		json
// @Produce		json
// @Param		person		body		handlers.InputPerson		true	"Person Object"
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/people	[POST]
func HandleCreatePerson(logger *httplog.Logger, service personCreator) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()

		// get and validate body as object
		personIn, problems, err := DecodeValidateBody[InputPerson](r)
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

		person, err:= service.CreatePerson(ctx, personIn)
		if err != nil {
			logger.Error("error creating object in database", "error", err)
			EncodeResponse(w, logger, http.StatusInternalServerError, ResponseErr{
				Error: "Error creating object",
			})
			return
		}

		// return response
		personOut := MapOutput(person)
		EncodeResponse(w, logger, http.StatusOK, ResponsePerson{
			Person: personOut,
		})
	}

}