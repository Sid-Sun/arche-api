package database

import (
	"database/sql"
	"go.uber.org/zap"
)

type DB struct {
	Users   UsersTable
	Folders FoldersTable
}

func NewDBInstance(dbClient *sql.DB, lgr *zap.Logger) *DB {
	return &DB{
		Users: &users{
			lgr: lgr,
			db:  dbClient,
		},
		Folders: &folders{
			lgr: lgr,
			db:  dbClient,
		},
	}
}
