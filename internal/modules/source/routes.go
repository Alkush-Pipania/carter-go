package source

import "github.com/go-chi/chi/v5"

func Routes(h *handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetSourcesByCollectionID)
	r.Post("/", h.CreateSource)
	return r
}
