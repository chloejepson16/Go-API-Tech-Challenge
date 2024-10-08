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

type courseLister interface{
	ListCourses(ctx context.Context) ([]models.Course, error)
	ListCourseByID(ctx context.Context, id int)(models.Course, error)
}

// HandleListCourses is a Handler that returns a list of all courses.
//
// @Summary		List all courses
// @Description	List all courses
// @Tags		courses
// @Accept		json
// @Produce		json
// @Success		200		{object}	handlers.responseMsg
// @Failure		500		{object}	handlers.responseErr
// @Router		/courses	[GET]
func HandleListCourses(logger *httplog.Logger, service courseLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		courses, err := service.ListCourses(ctx)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		coursesOut := mapMultipleCourseOutput(courses)
		encodeResponse(w, logger, http.StatusOK, responseCourses{
			Courses: coursesOut,
		})
	}
}

// HandleGetCourseById is a Handler that returns a single course by ID.
//
// @Summary     Get course by ID
// @Description Get a course by their ID
// @Tags        courses
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "id"
// @Success		200		{object}	handlers.responseMsg
// @Failure		500		{object}	handlers.responseErr
// @Router      /courses/{id}  [GET]
func HandleGetCourseById(logger *httplog.Logger, service courseLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()
		courseId := chi.URLParam(r, "id")
		fmt.Printf("ID is: %s", courseId)
		id, err := strconv.Atoi(courseId)
		fmt.Printf("ID 2 is: %d", id)

		// get values from database
		course, err := service.ListCourseByID(ctx, id)
		if err != nil {
			logger.Error("error getting specific people", "error", err)
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