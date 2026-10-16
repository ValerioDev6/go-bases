package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	httpresponse "web/internal/transport/http/response"
)

type TokenValidator interface {
	ValidateToken(token string) (jwt.MapClaims, error)
}

type authentication struct {
	claims jwt.MapClaims
}

type authenticationKey struct{}

func RequireAuthentication(validator TokenValidator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Fields(r.Header.Get("Authorization"))
		if len(authorization) != 2 || !strings.EqualFold(authorization[0], "Bearer") {
			httpresponse.Failure(w, http.StatusUnauthorized, httpresponse.GeneralError("bearer token is required"))
			return
		}

		claims, err := validator.ValidateToken(authorization[1])
		if err != nil {
			httpresponse.Failure(w, http.StatusUnauthorized, httpresponse.GeneralError("authentication token is invalid or expired"))
			return
		}

		value := authentication{claims: claims}
		ctx := context.WithValue(r.Context(), authenticationKey{}, value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticatedUserID(ctx context.Context) (int64, error) {
	authentication, _ := ctx.Value(authenticationKey{}).(authentication)
	subject, err := authentication.claims.GetSubject()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(subject, 10, 64)
}
