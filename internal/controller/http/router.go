package http

import (
	"net/http"

	"avito-task/internal/config"
	"avito-task/internal/controller/http/middleware"
	"avito-task/internal/controller/http/v1/auth"
	"avito-task/internal/controller/http/v1/merch"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	cfg config.Config,
	merchUsecase merch.Usecase,
	authUsecase auth.Usecase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Heartbeat("/health"))
	r.Use(chiMiddleware.Recoverer)

	r.Use(middleware.CheckJWT(cfg.JWTSecret, cfg.PublicRoutes))

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.CleanPath)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.Compress(4))
	r.Use(chiMiddleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", merch.NewHandler(merchUsecase).Router())
		r.Mount("/auth", auth.NewHandler(authUsecase).Router())
	})

	return r
}
