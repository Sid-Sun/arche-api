package service

import (
	"github.com/sid-sun/arche-api/app/database"
	"go.uber.org/zap"
)

type Service struct {
	Users   UsersService
	Folders FoldersService
}

func NewDBService(db *database.DB, lgr *zap.Logger) *Service {
	return &Service{
		Users: &users{
			db:  db,
			lgr: lgr,
		},
		Folders: &folders{
			db:  db,
			lgr: lgr,
		},
	}
}
