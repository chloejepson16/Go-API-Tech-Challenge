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
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/courses	[GET]
func HandleListCourses(logger *httplog.Logger, service courseLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		courses, err := service.ListCourses(ctx)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			handlers.EncodeResponse(w, logger, http.StatusInternalServerError, handlers.ResponseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		coursesOut := handlers.MapMultipleCourseOutput(courses)
		handlers.EncodeResponse(w, logger, http.StatusOK, handlers.ResponseCourses{
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
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
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
			handlers.EncodeResponse(w, logger, http.StatusInternalServerError, handlers.ResponseErr{
				Error: "Error retrieving data",
			})
			return
		}
		handlers.EncodeResponse(w, logger, http.StatusOK, handlers.ResponseCourse{
			Course: handlers.MapOutputCourse(course),
		})
	}
}