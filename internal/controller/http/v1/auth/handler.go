package auth

import "github.com/go-chi/chi/v5"

type Usecase interface{}

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", h.auth)

	return r
}
