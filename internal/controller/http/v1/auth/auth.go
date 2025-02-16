package auth

import (
	"net/http"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/go-chi/render"
)

const (
	authHandlerName = "auth"
)

func (h *Handler) auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.AuthUser

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, authHandlerName, "invalid request body", err)
		return
	}

	req.CurTime = time.Now().UTC()

	token, err := h.usecase.AuthUser(ctx, req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, authHandlerName, "invalid username or password", err)
		return
	}

	render.JSON(w, r, map[string]string{
		"token": token,
	})
}
