package database

import (
	"database/sql"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type UsersTable interface {
	Get(emailID string) (types.User, error)
	Create(emailID string, encryptionKey string, keyHash string) (types.UserID, error)
}

type users struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (u *users) Get(emailID string) (types.User, error) {
	query := `SELECT user_id, encryption_key, key_hash FROM users WHERE email=@email;`

	var userID types.UserID
	var encryptionKey, keyHash string

	row := u.db.QueryRow(query, sql.Named("email", emailID))
	err := row.Scan(&userID, &encryptionKey, &keyHash)
	if err != nil {
		return types.User{}, err
	}

	return types.User{
		ID:            userID,
		Email:         emailID,
		EncryptionKey: encryptionKey,
		KeyHash:       keyHash,
	}, nil
}

func (u *users) Create(emailID string, encryptionKey string, keyHash string) (types.UserID, error) {
	query := `INSERT INTO users (email, encryption_key, key_hash) OUTPUT inserted.user_id
VALUES (@email, @key, @hash);`

	var userID types.UserID

	row := u.db.QueryRow(query, sql.Named("email", emailID), sql.Named("key", encryptionKey), sql.Named("hash", keyHash))
	err := row.Err()
	if err != nil {
		return 0, err
	}
	err = row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
