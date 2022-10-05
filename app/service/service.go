package service

import (
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/initializers"
	"go.uber.org/zap"
)

type Service struct {
	Users   UsersService
	Folders FoldersService
	Notes   NotesService
}

func NewService(db *database.DB, mc initializers.MailClient, lgr *zap.Logger) *Service {
	return &Service{
		Users: &users{
			db:         db,
			lgr:        lgr,
			mailClient: mc,
		},
		Folders: &folders{
			db:  db,
			lgr: lgr,
		},
		Notes: &notes{
			db:  db,
			lgr: lgr,
		},
	}
}
