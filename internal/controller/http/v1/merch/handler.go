package merch

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Usecase interface {
	Info() error
	SendCoin() error
	BuyItem(item string) error
}

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

	r.Get("/info", h.info)
	r.Post("/send-coin", h.sendCoin)
	r.Post("/buy/{item}", h.buyItem)

	return r
}

func respondWithError(
	w http.ResponseWriter,
	statusCode int,
	handlerName string,
	errMsg string,
	err error,
) {
	slog.Error(err.Error(), "handler", handlerName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`{"error": "` + errMsg + `"}`))
}
