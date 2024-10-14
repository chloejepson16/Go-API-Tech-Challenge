package handlers

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
)

// HandleHelloWorld is a Handler that returns Hello, World... I am using it to set up scaffolding and test the applicaiton.
//
// @Summary		Hello, World!
// @Description A test endpoint to get everything up and running
// @Tags		test
// @Produce		json
// @Success		200		{object}	handlers.ResponseMsg
// @Failure		500		{object}	handlers.ResponseErr
// @Router		/	[GET]
func HandleHelloWorld(logger *httplog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		EncodeResponse(w, logger, http.StatusOK, ResponseMsg{
			Message: "Hello, World!",
		})
	}
}