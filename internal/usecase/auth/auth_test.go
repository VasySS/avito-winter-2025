package auth_test

import (
	"testing"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/dto"
	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecase_AuthUser(t *testing.T) {
	t.Parallel()

	const (
		hashCost = 10
	)

	type fields struct {
		repo      *mocks.Repository
		hasher    *mocks.PasswordHasher
		generator *mocks.TokenGenerator
	}

	tests := []struct {
		name    string
		fields  fields
		input   dto.AuthUser
		setup   func(f *fields, inp dto.AuthUser)
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "new user, successful registration",
			fields: fields{
				repo:      mocks.NewRepository(t),
				hasher:    mocks.NewPasswordHasher(t),
				generator: mocks.NewTokenGenerator(t),
			},
			input: dto.AuthUser{
				Username: "new_user",
				Password: "password",
				CurTime:  time.Now().UTC(),
			},
			setup: func(f *fields, inp dto.AuthUser) {
				f.repo.On("GetUserByUsername", mock.Anything, inp.Username).
					Return(entity.User{}, entity.ErrUserNotFound).Once()
				f.repo.On("CreateUser", mock.Anything, mock.AnythingOfType("entity.User")).
					Return(nil)
				f.repo.On("GetUserByUsername", mock.Anything, inp.Username).
					Return(entity.User{
						ID:        1,
						Username:  inp.Username,
						Password:  "hashed_password",
						Balance:   1000,
						CreatedAt: inp.CurTime,
					}, nil)

				f.hasher.On("GenerateFromPassword", []byte(inp.Password), hashCost).
					Return([]byte("hashed_password"), nil)
				f.generator.On("NewAccessToken", mock.AnythingOfType("entity.User"), inp.CurTime).
					Return("token123", nil)
			},
			want:    "token123",
			wantErr: assert.NoError,
		},
		{
			name: "existing user, correct password",
			fields: fields{
				repo:      mocks.NewRepository(t),
				hasher:    mocks.NewPasswordHasher(t),
				generator: mocks.NewTokenGenerator(t),
			},
			input: dto.AuthUser{
				Username: "existing_user",
				Password: "password123",
				CurTime:  time.Now().UTC(),
			},
			setup: func(f *fields, inp dto.AuthUser) {
				f.repo.On("GetUserByUsername", mock.Anything, inp.Username).
					Return(entity.User{Username: inp.Username, Password: "hashed_password", CreatedAt: inp.CurTime}, nil)
				f.hasher.On("CompareHashAndPassword", []byte("hashed_password"), []byte(inp.Password)).
					Return(nil)
				f.generator.On("NewAccessToken", mock.AnythingOfType("entity.User"), inp.CurTime).
					Return("token456", nil)
			},
			want:    "token456",
			wantErr: assert.NoError,
		},
		{
			name: "existing user, incorrect password",
			fields: fields{
				repo:      mocks.NewRepository(t),
				hasher:    mocks.NewPasswordHasher(t),
				generator: mocks.NewTokenGenerator(t),
			},
			input: dto.AuthUser{
				Username: "existing_user",
				Password: "wrong_password",
				CurTime:  time.Now().UTC(),
			},
			setup: func(f *fields, inp dto.AuthUser) {
				f.repo.On("GetUserByUsername", mock.Anything, inp.Username).
					Return(entity.User{Username: inp.Username, Password: "hashed_password"}, nil)
				f.hasher.On("CompareHashAndPassword", []byte("hashed_password"), []byte(inp.Password)).
					Return(assert.AnError)
			},
			want:    "",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := tt.fields
			uc := auth.New(
				fs.repo,
				fs.hasher,
				fs.generator,
			)

			tt.setup(&fs, tt.input)

			got, err := uc.AuthUser(t.Context(), tt.input)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
