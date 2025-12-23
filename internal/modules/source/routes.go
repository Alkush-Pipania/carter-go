package source

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetSourcesByCollectionID)
	r.Post("/", h.CreateSource)
	r.Post("/upload/presign", h.RequestUploadURL)
	r.Post("/{id}/confirm", h.ConfirmUpload)
	return r
}
