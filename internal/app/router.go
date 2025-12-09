package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(container *Container) http.Handler {
	//  1. Create the router (the Mux)
	r := chi.NewRouter()

	// 2. Middleware (global)
	// Logger: Logs every request to console (essential for debugging)
	r.Use(middleware.Logger)
	// Recoverer: Prevents the server from crashing if a handler panics (throws error)
	r.Use(middleware.Recoverer)

	// 3. Define the Health Route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// r.Route("/api/v1", func(v1Routes chi.Router){
	// 	v1Routes.Mount("/users", user.Routes())
	// })
	return r

}
