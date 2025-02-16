package auth_test

import (
	"testing"
	"time"

	"github.com/VasySS/avito-winter-2025/internal/entity"
	"github.com/VasySS/avito-winter-2025/internal/usecase/auth"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerator_NewAccessToken(t *testing.T) {
	t.Parallel()

	const (
		secret   = "mysecret"
		tokenTTL = 1 * time.Hour
	)

	user := entity.User{
		ID:       123,
		Username: "testuser",
	}

	currentTime := time.Now().UTC()

	t.Run("success", func(t *testing.T) {
		generator := auth.NewJWTGenerator(secret, tokenTTL)

		t.Parallel()

		token, err := generator.NewAccessToken(user, currentTime)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.ParseString(token, jwt.WithKey(jwa.HS256, []byte(secret)))
		require.NoError(t, err)

		expirationTime := parsedToken.Expiration()
		assert.WithinDuration(t, currentTime.Add(tokenTTL), expirationTime, 1*time.Second)

		id, ok := parsedToken.Get("id")
		require.True(t, ok)
		assert.EqualValues(t, user.ID, id)

		username, ok := parsedToken.Get("username")
		require.True(t, ok)
		assert.Equal(t, user.Username, username)
	})

	t.Run("missing secret", func(t *testing.T) {
		t.Parallel()

		generator := auth.NewJWTGenerator("", tokenTTL)

		_, err := generator.NewAccessToken(user, currentTime)
		require.Error(t, err)
	})
}
