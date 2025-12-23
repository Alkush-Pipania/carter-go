package app

import (
	"net/http"

	middl "github.com/Alkush-Pipania/carter-go/internal/middleware"
	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/source"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(container *Container) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middl.AuthMiddleware(container.jwt))
		r.Route("/api/v1", func(v1Routes chi.Router) {
			v1Routes.Mount("/users", user.Routes(container.userHandler))
			v1Routes.Mount("/collections", collection.Routes(container.collectionHandler))
			v1Routes.Mount("/sources", source.Routes(container.sourceHandler))
		})
	})

	return r
}
