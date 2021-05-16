package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func WithContentJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func ContextURLParams(lgr *zap.Logger, params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			m := make(map[string]string)
			for _, param := range params {
				m[param] = chi.URLParam(req, param)
			}
			ctx := context.WithValue(req.Context(), "url_params", m)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
