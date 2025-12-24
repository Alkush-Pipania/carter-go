package authentication

import "github.com/go-chi/chi/v5"

func Routes(h *handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.Post("/register", h.Register)

	return r
}
