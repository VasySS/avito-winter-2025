package merch

import (
	"net/http"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

const (
	sendCoinHandlerName = "send-coin"
)

func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) sendCoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, sendCoinHandlerName, "failed to get claims", err)
		return
	}

	senderUsername, ok := claims["username"].(string)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, sendCoinHandlerName, "failed to get sender username", nil)
		return
	}

	var req dto.CoinSend

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, "invalid request body", err)
		return
	}

	req.FromUser = senderUsername
	req.CurTime = time.Now().UTC()

	if err := h.usecase.SendCoin(ctx, req); err != nil {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, "failed to send coins", err)
		return
	}
}

func (h *Handler) buyItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, sendCoinHandlerName, "failed to get claims", err)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, sendCoinHandlerName, "failed to get username", nil)
		return
	}

	merchName := chi.URLParam(r, "item")
	if merchName == "" {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, "invalid item name", nil)
		return
	}

	req := dto.MerchPurchase{
		CurTime:   time.Now().UTC(),
		Username:  username,
		MerchName: merchName,
	}

	if err := h.usecase.BuyItem(ctx, req); err != nil {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, "failed to buy item", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
