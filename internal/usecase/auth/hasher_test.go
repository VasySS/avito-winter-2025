package auth_test

import (
	"testing"

	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordHasher_GenerateFromPassword(t *testing.T) {
	t.Parallel()

	type args struct {
		password []byte
		cost     int
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid password with default cost",
			args: args{
				password: []byte("secretpassword"),
				cost:     bcrypt.DefaultCost,
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty password",
			args: args{
				password: []byte(""),
				cost:     bcrypt.DefaultCost,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := auth.NewBcryptPasswordHasher()

			got, err := h.GenerateFromPassword(tt.args.password, tt.args.cost)
			tt.wantErr(t, err)

			assert.NotEmpty(t, got)
		})
	}
}

func TestBcryptPasswordHasher_CompareHashAndPassword(t *testing.T) {
	t.Parallel()

	password := []byte("secretpassword")

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	require.NoError(t, err)

	h := auth.NewBcryptPasswordHasher()

	err = h.CompareHashAndPassword(hashedPassword, password)
	require.NoError(t, err)

	err = h.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	require.Error(t, err)
}
