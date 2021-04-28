package service

import (
	"encoding/base64"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type UsersService interface {
	CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, error)
	GetUser(emailID string) (types.User, error)
}

type users struct {
	db  *database.DB
	lgr *zap.Logger
}

func (u *users) GetUser(emailID string) (types.User, error) {
	usr, err := u.db.Users.Get(emailID)
	if err != nil {
		// TODO: Add Logging
		// TODO: Check for no records error
		return types.User{}, err
	}

	return usr, nil
}

func (u *users) CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, error) {
	encryptionKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)
	hashStr := base64.StdEncoding.EncodeToString(keyHash[:])

	userID, err := u.db.Users.Create(emailID, encryptionKeyStr, hashStr)
	if err != nil {
		u.lgr.Sugar().Error(err)
		// TODO: Add Logging
		// TODO: Check for duplicate insertion error
		return types.User{}, err
	}

	return types.User{
		ID:            userID,
		Email:         emailID,
		KeyHash:       hashStr,
		EncryptionKey: encryptionKeyStr,
	}, nil
}
