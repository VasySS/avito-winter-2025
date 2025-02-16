package merch

import (
	"net/http"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	sendCoinHandlerName = "send-coin"
)

func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := getUsernameFromCtx(ctx, w)
	if username == "" {
		return
	}

	resp, err := h.usecase.Info(ctx, username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, "failed to get info", err)
		return
	}

	render.JSON(w, r, resp)
}

func (h *Handler) sendCoin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	senderUsername := getUsernameFromCtx(ctx, w)
	if senderUsername == "" {
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

	username := getUsernameFromCtx(ctx, w)
	if username == "" {
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
