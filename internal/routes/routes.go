package routes

import (
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/handlers"
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World"))
    })

	router.Get("/", handlers.HandleHelloWorld(logger))
	router.Get("/people", handlers.HandleListPeople(logger, svs))
	router.Get("/people/{id}", handlers.HandleGetPersonByID(logger, svs))
	router.Put("/people/{id}", handlers.HandleUpdatePerson(logger, svs))
	router.Delete("/people/{id}", handlers.HandleDeletePersonByID(logger, svs))
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
	router.Get("/courses", handlers.HandleListCourses(logger, svs))
}