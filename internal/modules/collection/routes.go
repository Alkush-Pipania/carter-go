package collection

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetCollections)     
	r.Post("/", h.CreateCollection)
	r.Put("/{id}", h.UpdateCollection)
	r.Delete("/{id}", h.DeleteCollection)
	return r
} 
