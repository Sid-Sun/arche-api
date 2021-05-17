package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type UsersTable interface {
	Get(emailID string) (types.User, *erx.Erx)
	Create(emailID string, encryptionKey string, keyHash string) (types.UserID, *erx.Erx)
}

type users struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (u *users) Get(emailID string) (types.User, *erx.Erx) {
	query := `SELECT user_id, encryption_key, key_hash FROM users WHERE email=@email;`

	var userID types.UserID
	var encryptionKey, keyHash string

	row := u.db.QueryRow(query, sql.Named("email", emailID))
	err := row.Scan(&userID, &encryptionKey, &keyHash)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [Get] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return types.User{}, errx
		}
		if errors.Is(err, sql.ErrNoRows) {
			errx = erx.WithArgs(errx, erx.SeverityInfo, custom_errors.NoRowsInResultSet)
			u.lgr.Info(fmt.Sprintf("[Database] [Users] [Get] [Scan] [ErrSQLNoResultsInSet] %s", errx.String()))
			return types.User{}, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Get] [Scan] %s", errx.Error()))
		return types.User{}, errx
	}

	return types.User{
		ID:            userID,
		Email:         emailID,
		EncryptionKey: encryptionKey,
		KeyHash:       keyHash,
	}, nil
}

func (u *users) Create(emailID string, encryptionKey string, keyHash string) (types.UserID, *erx.Erx) {
	query := `INSERT INTO users (email, encryption_key, key_hash) OUTPUT inserted.user_id
VALUES (@email, @key, @hash);`

	var userID types.UserID

	row := u.db.QueryRow(query, sql.Named("email", emailID), sql.Named("key", encryptionKey), sql.Named("hash", keyHash))
	err := row.Err()
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [Create] [QueryRow] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Create] [QueryRow] %s", err.Error()))
		return 0, errx
	}

	err = row.Scan(&userID)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			switch sqlErr.Number {
			case 2601:
				errx = erx.WithArgs(errx, custom_errors.DuplicateRecordInsertion, erx.SeverityInfo)
				u.lgr.Info(fmt.Sprintf("[Database] [Users] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			default:
				errx = erx.WithArgs(errx, erx.SeverityError)
				u.lgr.Error(fmt.Sprintf("[Database] [Users] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			}
			return 0, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Create] [Scan] %s", errx.Error()))
		return 0, errx
	}

	return userID, nil
}
