package router

import (
	"github.com/go-chi/chi"
	"github.com/sid-sun/arche-api/app/handlers"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func NewRouter(svc *service.Service, jwtCfg *config.JWTConfig, lgr *zap.Logger) *chi.Mux {
	rtr := chi.NewRouter()

	rtr.Post("/v1/signup", handlers.CreateUserHandler(svc.Users, jwtCfg, lgr))
	rtr.Post("/v1/login", handlers.LoginUserHandler(svc.Users, jwtCfg, lgr))

	rtr.Post("/v1/folders/create", handlers.CreateFolderHandler(svc.Folders, jwtCfg, lgr))
	rtr.Get("/v1/folders/get", handlers.GetFoldersHandler(svc.Folders, jwtCfg, lgr))
	rtr.Delete("/v1/folders/delete", handlers.DeleteFolder(svc.Folders, jwtCfg, lgr))

	return rtr
}
