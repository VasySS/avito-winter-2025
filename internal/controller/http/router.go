package http

import (
	"net/http"

	"github.com/VasySS/avito-winter-2025/internal/config"
	"github.com/VasySS/avito-winter-2025/internal/controller/http/middleware"
	"github.com/VasySS/avito-winter-2025/internal/controller/http/v1/auth"
	"github.com/VasySS/avito-winter-2025/internal/controller/http/v1/merch"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	cfg config.Config,
	merchUsecase merch.Usecase,
	authUsecase auth.Usecase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
		chiMiddleware.Heartbeat("/health"),
		middleware.CheckJWT(cfg.JWTSecret, cfg.PublicRoutes),
		chiMiddleware.RequestID,
		chiMiddleware.CleanPath,
		chiMiddleware.StripSlashes,
		chiMiddleware.Compress(5),
	)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", merch.NewHandler(merchUsecase).Router())
		r.Mount("/auth", auth.NewHandler(authUsecase).Router())
	})

	return r
}
