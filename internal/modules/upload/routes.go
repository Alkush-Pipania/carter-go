package upload

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/presign", h.RequestUploadURL)
	r.Post("/{id}/confirm", h.ConfirmUpload)
	return r
}
