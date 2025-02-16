package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VasySS/avito-winter-2025/internal/config"
	httpRouter "github.com/VasySS/avito-winter-2025/internal/controller/http"
	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/repository/postgres"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

const (
	migrationsPath = "../../migrations"
	envPath        = "../../.env"
	authPath       = "/api/auth"
)

type HandlerTestSuite struct {
	suite.Suite
	pgPool   *pgxpool.Pool
	pgFacade *postgres.Facade
	router   http.Handler
}

func (s *HandlerTestSuite) SetupSuite() {
	ctx := context.Background()
	cfg := config.MustInit(envPath)

	gofakeit.GlobalFaker = gofakeit.NewFaker(source.NewCrypto(), false)

	connURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)

	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		s.T().Fatalf("failed to connect to database: %v", err)
	}

	txManager := postgres.NewTxManager(pool)
	pgStorage := postgres.New(txManager)
	pgFacade := postgres.NewFacade(pgStorage)

	authUsecase := auth.New(
		pgFacade,
		auth.NewBcryptPasswordHasher(),
		auth.NewJWTGenerator(cfg.JWTSecret, cfg.AccessTokenTTL),
	)

	s.pgPool = pool
	s.pgFacade = pgFacade
	s.router = httpRouter.NewRouter(cfg, nil, authUsecase)
}

func (s *HandlerTestSuite) TearDownSuite() {
	s.pgPool.Close()
}

func TestAuthHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping merch integration test in short mode.")
	}

	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestAuthHandler() {
	randUsername := gofakeit.Username()

	reqBody, err := json.Marshal(dto.AuthUser{
		Username: randUsername,
		Password: gofakeit.Password(true, true, true, true, false, 10),
	})
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPost, authPath, bytes.NewBuffer(reqBody))
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	s.Equal(http.StatusOK, rr.Code)

	var respBody struct {
		Token string `json:"token"`
	}

	err = json.NewDecoder(rr.Body).Decode(&respBody)
	s.Require().NoError(err)
	s.NotEmpty(respBody.Token)

	badReqBody, err := json.Marshal(dto.AuthUser{
		Username: randUsername,
		Password: "wrongpass",
	})
	s.NoError(err)

	badReq, err := http.NewRequest(http.MethodPost, authPath, bytes.NewBuffer(badReqBody))
	s.Require().NoError(err)

	badRR := httptest.NewRecorder()
	s.router.ServeHTTP(badRR, badReq)
	s.Equal(http.StatusUnauthorized, badRR.Code)
}
