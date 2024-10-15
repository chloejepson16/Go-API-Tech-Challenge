package routes

import (
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
	courseHandlers "github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers/courses"
	peopleHandlers "github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers/people"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type Option func (*routerOptions)

type routerOptions struct{
	registerHealthRoute bool
}

func WithRegisterHealthRoute(registerHealthRoute bool) Option{
	return func(options *routerOptions){
		options.registerHealthRoute= registerHealthRoute
	}
}

func RegisterRoutes(router *chi.Mux, logger *httplog.Logger, svs *services.PersonService, opts ...Option){
	options:= routerOptions{
		registerHealthRoute: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.registerHealthRoute {
		router.Get("/health-check", handlers.HandleHealth(logger))
	}

	router.Get("/", handlers.HandleHelloWorld(logger))
	router.Get("/people", peopleHandlers.HandleListPeople(logger, svs))
	router.Get("/people/{id}", peopleHandlers.HandleGetPersonByID(logger, svs))
	router.Put("/people/{id}", peopleHandlers.HandleUpdatePerson(logger, svs))
	router.Delete("/people/{id}", peopleHandlers.HandleDeletePersonByID(logger, svs))
	router.Post("/people", peopleHandlers.HandleCreatePerson(logger, svs))
}

func RegisterRoutesForCourses(router *chi.Mux, logger *httplog.Logger, svs *services.CourseService, opts ...Option){
	options:= routerOptions{
		registerHealthRoute: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.registerHealthRoute {
		router.Get("/health-check", handlers.HandleHealth(logger))
	}
	router.Get("/courses", courseHandlers.HandleListCourses(logger, svs))
	router.Get("/courses/{id}", courseHandlers.HandleGetCourseById(logger, svs))
	router.Put("/courses/{id}", courseHandlers.HandleUpdateCourse(logger, svs))
	router.Delete("/courses/{id}", courseHandlers.HandleDeleteCourseByID(logger, svs))
	router.Post("/courses", courseHandlers.HandleCreateCourse(logger, svs))
}