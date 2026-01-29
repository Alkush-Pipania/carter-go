package app

import (
	"net/http"

	middl "github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/internal/modules/authentication"
	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/source"
	"github.com/Alkush-Pipania/carter-go/internal/modules/upload"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(container *Container) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Turnstile-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public auth routes (no authentication required)
	r.Mount("/api/v1/auth", authentication.Routes(container.authHandler))

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(middl.AuthMiddleware(container.authService))
		r.Route("/api/v1", func(v1Routes chi.Router) {
			v1Routes.Mount("/users", user.Routes(container.userHandler))
			v1Routes.Mount("/collections", collection.Routes(container.collectionHandler))
			v1Routes.Mount("/sources", source.Routes(container.sourceHandler))
			v1Routes.Mount("/upload", upload.Routes(container.uploadHandler))
		})
	})

	return r
}
