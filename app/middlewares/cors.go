package middlewares

import (
	"github.com/rs/cors"
	"net/http"
)

func WithCors() func(h http.Handler) http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Refresh_Token"},
		MaxAge:         30 * 60, // 30 mins of preflight caching
	}).Handler

	return handler
}
