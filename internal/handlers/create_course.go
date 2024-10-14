package handlers

import (
	"context"
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
)

type courseCreator interface{
	CreateCourse(ctx context.Context, course models.Course) (models.Course, error)
}

// HandleCreateCourse is a Handler that updates a user based on a user object from the request body.
//
// @Summary		Create a course
// @Description	Create a course
// @Tags		courses
// @Accept		json
// @Produce		json
// @Param		course		body		handlers.InputCourse		true	"Course Object"
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/courses	[POST]
func HandleCreateCourse(logger *httplog.Logger, service courseCreator) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()

		// get and validate body as object
		courseIn, problems, err := DecodeValidateBody[InputCourse](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				EncodeResponse(w, logger, http.StatusBadRequest, ResponseErr{
					Error: "error validating courseIn",
				})
			default:
				logger.Error("BodyParser error", "error", err)
				EncodeResponse(w, logger, http.StatusBadRequest, ResponseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		course, err:= service.CreateCourse(ctx, courseIn)
		if err != nil {
			logger.Error("error creating object in database", "error", err)
			EncodeResponse(w, logger, http.StatusInternalServerError, ResponseErr{
				Error: "Error creating object",
			})
			return
		}

		// return response
		courseOut := MapOutputCourse(course)
		EncodeResponse(w, logger, http.StatusOK, ResponseCourse{
			Course: courseOut,
		})
	}

}