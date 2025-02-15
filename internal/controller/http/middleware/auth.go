package middleware

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func CheckJWT(jwtSecret string, publicRoutes []string) func(next http.Handler) http.Handler {
	tokenAuth := jwtauth.New(
		string(jwa.HS256),
		[]byte(jwtSecret),
		nil,
		jwt.WithAcceptableSkew(10*time.Second),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(publicRoutes, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			token, err := jwtauth.VerifyRequest(tokenAuth, r, tokenFromHeader, tokenFromCookie)
			if err != nil {
				slog.Debug(err.Error())
				http.Error(w, "error authorizing user", http.StatusUnauthorized)

				return
			}

			ctx := jwtauth.NewContext(r.Context(), token, nil)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// tokenFromCookie tries to retrieve the token string from a cookie named
// "accessToken".
func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("accessToken")
	if err != nil {
		return ""
	}

	return cookie.Value
}

// tokenFromHeader tries to retrieve the token string from the
// "Authorization" request header: "Authorization: BEARER T".
func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	return ""
}
