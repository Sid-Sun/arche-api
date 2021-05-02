package middlewares

import (
	"context"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"net/http"
)

func JWTAuth(jwtCfg *config.JWTConfig, lgr *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("authentication_token")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

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
