package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type courseUpdater interface {
	UpdateCourse(ctx context.Context, id int, course models.Course) (models.Course, error)
}

// HandleUpdateCourse is a Handler that updates a user based on a user object from the request body.
//
// @Summary		Update a course by id
// @Description	Update a course by id
// @Tags		courses
// @Accept		json
// @Produce		json
// @Param		id			path		int	true						"id"
// @Param		course		body		handlers.inputCourse		true	"Course Object"
// @Success		200		{object}	handlers.responseMsg
// @Failure		500		{object}	handlers.responseErr
// @Router		/courses/{id}	[PUT]
func HandleUpdateCourse(logger *httplog.Logger, service courseUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get and validate ID
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get and validate body as object
		courseIn, problems, err := decodeValidateBody[inputCourse, models.Course](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "error validating courseIn",
				})
			default:
				logger.Error("BodyParser error", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		// update object in database
		course, err := service.UpdateCourse(ctx, id, courseIn)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating object",
			})
			return
		}

		// return response
		courseOut := mapOutputCourse(course)
		encodeResponse(w, logger, http.StatusOK, responseCourse{
			Course: courseOut,
		})
	}
}