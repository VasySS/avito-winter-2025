package merch_test

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
	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/VasySS/avito-winter-2025/internal/repository/postgres"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"
	merchUc "github.com/VasySS/avito-winter-2025/internal/usecase/merch"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

const (
	migrationsPath = "../../migrations"
	envPath        = "../../.env"
	authPath       = "/api/auth"
	buyPath        = "/api/buy"
	infoPath       = "/api/info"
	sendCoinPath   = "/api/send-coin"
)

type HandlerTestSuite struct {
	suite.Suite
	pgPool   *pgxpool.Pool
	pgFacade *postgres.Facade
	router   http.Handler
	tokenStr string
	username string
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

	merchUsecase := merchUc.New(pgFacade)
	authUsecase := auth.New(
		pgFacade,
		auth.NewBcryptPasswordHasher(),
		auth.NewJWTGenerator(cfg.JWTSecret, cfg.AccessTokenTTL),
	)

	s.pgPool = pool
	s.pgFacade = pgFacade
	s.router = httpRouter.NewRouter(cfg, merchUsecase, authUsecase)
}

func (s *HandlerTestSuite) TearDownSuite() {
	s.pgPool.Close()
}

func (s *HandlerTestSuite) SetupTest() {
	s.userLogin()
}

func (s *HandlerTestSuite) userLogin() {
	randUsername := gofakeit.Username()

	reqBody, err := json.Marshal(dto.AuthUser{
		Username: randUsername,
		Password: gofakeit.Password(true, true, true, true, false, 10),
	})
	s.Require().NoError(err)

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

	s.tokenStr = respBody.Token
	s.username = randUsername
}

func TestMerchHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping merch integration test in short mode.")
	}

	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) TestBuyItemHandler() {
	const itemName = "t-shirt"

	req, err := http.NewRequest(http.MethodPost, buyPath+"/"+itemName, nil)
	req.Header.Set("Authorization", "Bearer "+s.tokenStr)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)

	user, err := s.pgFacade.GetUserByUsername(s.T().Context(), s.username)
	s.NoError(err)
	s.Equal(920, user.Balance)

	resp, err := s.pgFacade.Info(s.T().Context(), user.ID)
	s.NoError(err)
	s.Equal(1, resp.Inventory[0].Quantity)
	s.Equal(itemName, resp.Inventory[0].Name)
}

func (s *HandlerTestSuite) TestSendCoinHandler() {
	secondUser := entity.User{
		Username:  gofakeit.Username(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		CreatedAt: gofakeit.PastDate(),
	}

	s.NoError(s.pgFacade.CreateUser(s.T().Context(), secondUser))

	reqBody, err := json.Marshal(dto.CoinSend{
		ToUser: secondUser.Username,
		Amount: 123,
	})
	s.NoError(err)

	req, err := http.NewRequest(http.MethodPost, sendCoinPath, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+s.tokenStr)
	s.NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)

	sender, err := s.pgFacade.GetUserByUsername(s.T().Context(), s.username)
	s.NoError(err)
	s.Equal(877, sender.Balance)

	receiver, err := s.pgFacade.GetUserByUsername(s.T().Context(), secondUser.Username)
	s.NoError(err)
	s.Equal(1123, receiver.Balance)
}

func (s *HandlerTestSuite) TestInfoHandler() {
	secondUser := entity.User{
		Username:  gofakeit.Username(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		CreatedAt: gofakeit.PastDate(),
	}

	s.NoError(s.pgFacade.CreateUser(s.T().Context(), secondUser))

	firstUserDB, err := s.pgFacade.GetUserByUsername(s.T().Context(), s.username)
	s.NoError(err)

	secondUserDB, err := s.pgFacade.GetUserByUsername(s.T().Context(), secondUser.Username)
	s.NoError(err)

	merch, err := s.pgFacade.GetMerch(s.T().Context(), "t-shirt")
	s.NoError(err)

	s.NoError(s.pgFacade.SendCoins(s.T().Context(), entity.UserTransfer{
		SenderUserID:   firstUserDB.ID,
		ReceiverUserID: secondUserDB.ID,
		Amount:         merch.Price,
		CreatedAt:      gofakeit.PastDate(),
	}))

	s.NoError(s.pgFacade.BuyMerch(s.T().Context(), entity.MerchPurchase{
		UserID:      firstUserDB.ID,
		MerchItemID: merch.ID,
		Price:       80,
		CreatedAt:   gofakeit.PastDate(),
	}))

	req, err := http.NewRequest(http.MethodGet, infoPath, nil)
	req.Header.Set("Authorization", "Bearer "+s.tokenStr)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)

	var resp dto.InfoResponse

	err = json.NewDecoder(rr.Body).Decode(&resp)
	s.Require().NoError(err)
	s.Equal(1000-merch.Price-80, resp.Coins)
	s.Equal(1, resp.Inventory[0].Quantity)
	s.Equal(merch.Name, resp.Inventory[0].Name)
}
