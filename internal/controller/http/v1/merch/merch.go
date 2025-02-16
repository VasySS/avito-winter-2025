package merch

import (
	"errors"
	"net/http"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	sendCoinHandlerName = "send-coin"
	buyItemHandlerName  = "buy-item"
	infoHandlerName     = "info"
)

func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := getUsernameFromCtx(ctx, w)
	if username == "" {
		return
	}

	resp, err := h.usecase.Info(ctx, username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, infoHandlerName, "failed to get info", err)
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

	err := h.usecase.SendCoin(ctx, req)
	if errors.Is(err, entity.ErrLowBalance) {
		respondWithError(w, http.StatusBadRequest, sendCoinHandlerName, err.Error(), err)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, sendCoinHandlerName, "failed to send coins", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) buyItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := getUsernameFromCtx(ctx, w)
	if username == "" {
		return
	}

	merchName := chi.URLParam(r, "item")
	if merchName == "" {
		respondWithError(w, http.StatusBadRequest, buyItemHandlerName, "invalid item name", nil)
		return
	}

	req := dto.MerchPurchase{
		CurTime:   time.Now().UTC(),
		Username:  username,
		MerchName: merchName,
	}

	err := h.usecase.BuyItem(ctx, req)
	if errors.Is(err, entity.ErrLowBalance) {
		respondWithError(w, http.StatusBadRequest, buyItemHandlerName, err.Error(), err)
		return
	}

	if errors.Is(err, entity.ErrSameTransferUser) {
		respondWithError(w, http.StatusBadRequest, buyItemHandlerName, err.Error(), err)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, buyItemHandlerName, "failed to buy item", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
