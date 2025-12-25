package source

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateSource)
	r.Get("/collection/{id}", h.GetSourcesByCollectionID)
	r.Get("/{id}", h.GetSourceByID)
	r.Delete("/{id}", h.DeleteSource)
	return r
}
