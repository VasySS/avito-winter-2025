package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/go-chi/chi/v5"
)

type Usecase interface {
	AuthUser(ctx context.Context, req dto.AuthUser) (string, error)
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

	r.Post("/", h.auth)

	return r
}

func respondWithError(
	w http.ResponseWriter,
	statusCode int,
	handlerName string,
	errMsg string,
	err error,
) {
	slog.Error(err.Error(), "auth-handler", handlerName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`{"error": "` + errMsg + `"}`))
}
