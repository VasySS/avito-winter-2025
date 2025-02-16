package merch_test

import (
	"testing"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/VasySS/avito-winter-2025/internal/usecase/merch"
	"github.com/VasySS/avito-winter-2025/internal/usecase/merch/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_SendCoin(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo *mocks.Repository
	}

	tests := []struct {
		name    string
		fields  fields
		input   dto.CoinSend
		setup   func(f *fields, inp dto.CoinSend)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successful send coin",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.CoinSend{
				FromUser: "user1",
				ToUser:   "user2",
				Amount:   100,
			},
			setup: func(f *fields, inp dto.CoinSend) {
				f.repo.On("GetUserByUsername", t.Context(), inp.FromUser).
					Return(entity.User{ID: 1, Username: inp.FromUser, Balance: 1000}, nil).Once()
				f.repo.On("GetUserByUsername", t.Context(), inp.ToUser).
					Return(entity.User{ID: 2, Username: inp.ToUser, Balance: 500}, nil).Once()

				f.repo.On("SendCoins", t.Context(), entity.UserTransfer{
					SenderUserID:   1,
					ReceiverUserID: 2,
					Amount:         inp.Amount,
					CreatedAt:      inp.CurTime,
				}).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "receiver not found",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.CoinSend{
				FromUser: "user1",
				ToUser:   "user2",
				Amount:   100,
			},
			setup: func(f *fields, inp dto.CoinSend) {
				f.repo.On("GetUserByUsername", t.Context(), inp.FromUser).
					Return(entity.User{ID: 1, Username: inp.FromUser, Balance: 1000}, nil).Once()
				f.repo.On("GetUserByUsername", t.Context(), inp.ToUser).
					Return(entity.User{}, entity.ErrUserNotFound).Once()
			},
			wantErr: assert.Error,
		},
		{
			name: "receiver is the same as sender",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.CoinSend{
				FromUser: "user1",
				ToUser:   "user1",
				Amount:   100,
			},
			setup:   func(_ *fields, _ dto.CoinSend) {},
			wantErr: assert.Error,
		},
		{
			name: "not enough coins",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.CoinSend{
				FromUser: "user1",
				ToUser:   "user2",
				Amount:   1000,
			},
			setup: func(f *fields, inp dto.CoinSend) {
				f.repo.On("GetUserByUsername", t.Context(), inp.FromUser).
					Return(entity.User{ID: 1, Username: inp.FromUser, Balance: 999}, nil).Once()
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tt.fields
			uc := merch.New(
				fs.repo,
			)

			tt.setup(&fs, tt.input)

			err := uc.SendCoin(t.Context(), tt.input)
			tt.wantErr(t, err)
		})
	}
}

func TestUsecase_BuyItem(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo *mocks.Repository
	}

	tests := []struct {
		name    string
		fields  fields
		input   dto.MerchPurchase
		setup   func(f *fields, inp dto.MerchPurchase)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successful purchase",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.MerchPurchase{
				Username:  "user1",
				MerchName: "merch1",
			},
			setup: func(f *fields, inp dto.MerchPurchase) {
				const (
					merchID    = 1
					merchPrice = 100
					userID     = 1
				)

				f.repo.On("GetUserByUsername", t.Context(), inp.Username).
					Return(entity.User{ID: userID, Username: inp.Username, Balance: 1000}, nil).Once()
				f.repo.On("GetMerch", t.Context(), inp.MerchName).
					Return(entity.Merch{ID: merchID, Name: inp.MerchName, Price: merchPrice}, nil).Once()

				f.repo.On("BuyMerch", t.Context(), entity.MerchPurchase{
					UserID:      userID,
					MerchItemID: merchID,
					Price:       merchPrice,
					CreatedAt:   inp.CurTime,
				}).Return(nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "balance not enough",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.MerchPurchase{
				Username:  "user1",
				MerchName: "merch1",
			},
			setup: func(f *fields, inp dto.MerchPurchase) {
				f.repo.On("GetUserByUsername", t.Context(), inp.Username).
					Return(entity.User{ID: 1, Username: inp.Username, Balance: 99}, nil).Once()
				f.repo.On("GetMerch", t.Context(), inp.MerchName).
					Return(entity.Merch{ID: 1, Name: inp.MerchName, Price: 100}, nil).Once()
			},
			wantErr: assert.Error,
		},
		{
			name: "merch not found",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: dto.MerchPurchase{
				Username:  "user1",
				MerchName: "merch_not_existing",
			},
			setup: func(f *fields, inp dto.MerchPurchase) {
				f.repo.On("GetMerch", t.Context(), inp.MerchName).
					Return(entity.Merch{}, entity.ErrMerchItemNotFound).Once()
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tt.fields
			u := merch.New(
				fs.repo,
			)

			tt.setup(&fs, tt.input)

			err := u.BuyItem(t.Context(), tt.input)
			tt.wantErr(t, err)
		})
	}
}

func TestUsecase_Info(t *testing.T) {
	t.Parallel()

	resp := dto.InfoResponse{
		Coins: 1000,
		Inventory: []dto.InventoryItem{
			{Name: "item1", Quantity: 1},
		},
		CoinHistory: dto.CoinHistory{
			Received: []dto.CoinTransferReceived{
				{FromUser: "user2", Amount: 100},
			},
			Sent: []dto.CoinTransferSent{
				{ToUser: "user2", Amount: 100},
			},
		},
	}

	type fields struct {
		repo *mocks.Repository
	}

	tests := []struct {
		name    string
		fields  fields
		input   string
		setup   func(f *fields, inp string)
		want    dto.InfoResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successful get",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: "user1",
			setup: func(f *fields, inp string) {
				f.repo.On("GetUserByUsername", t.Context(), inp).
					Return(entity.User{ID: 1, Username: inp}, nil).Once()
				f.repo.On("Info", t.Context(), int64(1)).
					Return(resp, nil).Once()
			},
			want:    resp,
			wantErr: assert.NoError,
		},
		{
			name: "unsuccessful get",
			fields: fields{
				repo: mocks.NewRepository(t),
			},
			input: "user1",
			setup: func(f *fields, inp string) {
				f.repo.On("GetUserByUsername", t.Context(), inp).
					Return(entity.User{ID: 1, Username: inp}, nil).Once()
				f.repo.On("Info", t.Context(), int64(1)).
					Return(dto.InfoResponse{}, assert.AnError).Once()
			},
			want:    dto.InfoResponse{},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tt.fields
			uc := merch.New(
				fs.repo,
			)

			tt.setup(&fs, tt.input)

			got, err := uc.Info(t.Context(), tt.input)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
