package middlewares

import (
	"context"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func JWTAuth(jwtCfg *config.JWTConfig, lgr *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Otherwise, extract the actual token
			hasToken := strings.Split(token, "Bearer")
			if len(hasToken) != 2 || len(hasToken[1]) <= 1 { // Compensating for a space literal
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			token = strings.TrimSpace(hasToken[1])

			var claims types.AccessTokenClaims
			var err error
			if claims, err = utils.ValidateJWT(token, jwtCfg.GetSecret(), lgr); err != nil {
				// TODO: Add Logging
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			// just a stub.. some ideas are to look at URL query params for something like
			// the page number, or the limit, and send a query cursor down the chain
			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), "claims", claims)))
		})
	}
}
