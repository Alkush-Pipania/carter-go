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
)

func NewRouter(
	userHandler *user.UserHandler,
	collectionHandler *collection.Handler,
	sourceHandler *source.Handler,
	uploadHandler *upload.Handler,
	authHandler *authentication.Handler,
	authService *authentication.Service,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public auth routes (no authentication required)
	r.Mount("/api/v1/auth", authentication.Routes(authHandler))

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(middl.AuthMiddleware(authService))
		r.Route("/api/v1", func(v1Routes chi.Router) {
			v1Routes.Mount("/users", user.Routes(userHandler))
			v1Routes.Mount("/collections", collection.Routes(collectionHandler))
			v1Routes.Mount("/sources", source.Routes(sourceHandler))
			v1Routes.Mount("/upload", upload.Routes(uploadHandler))
		})
	})

	return r
}
