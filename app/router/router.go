package router

import (
	"github.com/go-chi/chi"
	"github.com/sid-sun/arche-api/app/handlers"
	"github.com/sid-sun/arche-api/app/middlewares"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func NewRouter(svc *service.Service, jwtCfg *config.JWTConfig, lgr *zap.Logger) *chi.Mux {
	rtr := chi.NewRouter()

	rtr.Post("/v1/signup", handlers.CreateUserHandler(svc.Users, jwtCfg, lgr))
	rtr.Post("/v1/login", handlers.LoginUserHandler(svc.Users, jwtCfg, lgr))

	rtr.Route("/v1/folders", func(r chi.Router) {
		r.Use(middlewares.JWTAuth(jwtCfg, lgr))

		r.Post("/create", handlers.CreateFolderHandler(svc.Folders, jwtCfg, lgr))
		r.Get("/get", handlers.GetFoldersHandler(svc.Folders, jwtCfg, lgr))
		r.Delete("/delete", handlers.DeleteFolder(svc.Folders, jwtCfg, lgr))
	})

	return rtr
}
