package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sid-sun/arche-api/app/handlers"
	"github.com/sid-sun/arche-api/app/middlewares"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func NewRouter(svc *service.Service, jwtCfg *config.JWTConfig, lgr *zap.Logger) *chi.Mux {
	rtr := chi.NewRouter()

	rtr.Use(middleware.Recoverer)
	rtr.Use(middlewares.WithContentJSON)

	rtr.Post("/v1/signup", handlers.CreateUserHandler(svc.Users, jwtCfg, lgr))
	rtr.Post("/v1/login", handlers.LoginUserHandler(svc.Users, jwtCfg, lgr))

	rtr.Route("/v1/session", func(r chi.Router) {
		r.With(middlewares.JWTAuth(jwtCfg, lgr)).Get("/validate", handlers.ValidateTokenHandler(lgr))
		r.Post("/refresh", handlers.RefreshTokenHandler(jwtCfg, lgr))
	})

	rtr.Route("/v1/folders", func(r chi.Router) {
		r.Use(middlewares.JWTAuth(jwtCfg, lgr))

		r.Post("/create", handlers.CreateFolderHandler(svc.Folders, lgr))
		r.Get("/get", handlers.GetFoldersHandler(svc.Folders, lgr))
		r.Delete("/delete", handlers.DeleteFolderHandler(svc.Folders, lgr))
	})

	rtr.Route("/v1/notes", func(r chi.Router) {
		r.Use(middlewares.JWTAuth(jwtCfg, lgr))

		r.Post("/create", handlers.CreateNoteHandler(svc.Notes, lgr))
		r.Put("/update", handlers.UpdateNoteHandler(svc.Notes, lgr))
		r.Get("/getall", handlers.GetNotesHandler(svc.Notes, lgr))
		r.Delete("/delete", handlers.DeleteNoteHandler(svc.Notes, lgr))
	})

	return rtr
}
