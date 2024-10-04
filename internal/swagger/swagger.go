package swagger

import (
	"fmt"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/swagger/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RunSwagger(r *chi.Mux, logger *httplog.Logger, host string) {
	// docs
	docs.SwaggerInfo.Title = "User Microservice API"
	docs.SwaggerInfo.Description = "Sample Go API"
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/"

	docs.SwaggerInfo.Schemes = []string{"http"}

	// handler
	baseURL := "http://" + host

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(baseURL+"/swagger/doc.json"),
	))

	logger.Info(fmt.Sprintf("Swagger URL: %s/swagger/index.html", baseURL))
}