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
			setup:   func(f *fields, inp dto.CoinSend) {},
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
