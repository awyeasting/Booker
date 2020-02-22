package queue

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"net/http"
)

// Route all book club queue related API handles to the proper handlers
func QueueRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))

	// TODO: Actually do routing

	return r
}