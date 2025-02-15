package merch

import "github.com/go-chi/chi/v5"

type Usecase interface {
}

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/info", h.info)
	r.Post("/send-coin", h.sendCoin)
	r.Post("/buy/{item}", h.buyItem)

	return r
}
