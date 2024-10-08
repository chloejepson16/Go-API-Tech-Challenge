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

type courseDeleter interface{
	DeleteCourseByID(ctx context.Context, id int) (models.Course, error)
}

// HandleDeleteCourseByID is a Handler that deletes a single course by ID.
//
// @Summary     Delete course by ID
// @Description Delete a course by their ID
// @Tags        courses
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "id"
// @Success     200  {object}    handlers.responseMsg
// @Failure     500  {object}    handlers.responseErr
// @Router      /courses/{id}  [DELETE]
func HandleDeleteCourseByID(logger *httplog.Logger, service courseDeleter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		// setup
		ctx := r.Context()
		courseId := chi.URLParam(r, "id")
		fmt.Printf("ID is: %s", courseId)
		id, err := strconv.Atoi(courseId)
		fmt.Printf("ID 2 is: %d", id)

		course, err:= service.DeleteCourseByID(ctx, id)
		if err != nil{
			logger.Error("error deleting specific people", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}
		encodeResponse(w, logger, http.StatusOK, responseCourse{
			Course: mapOutputCourse(course),
		})
	}
}