package service

import (
	"encoding/base64"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type UsersService interface {
	CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, *erx.Erx)
	GetUser(emailID string) (types.User, *erx.Erx)
}

type users struct {
	db  *database.DB
	lgr *zap.Logger
}

func (u *users) GetUser(emailID string) (types.User, *erx.Erx) {
	usr, errx := u.db.Users.Get(emailID)
	if errx != nil {
		(*u).lgr.Debug(fmt.Sprintf("[Service] [Users] [GetUser] [Get] %s", errx.Error()))
		return types.User{}, errx
	}

	return usr, nil
}

func (u *users) CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, *erx.Erx) {
	encryptionKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)
	hashStr := base64.StdEncoding.EncodeToString(keyHash[:])

	userID, errx := u.db.Users.Create(emailID, encryptionKeyStr, hashStr)
	if errx != nil {
		(*u).lgr.Debug(fmt.Sprintf("[Service] [Users] [CreateUser] [Create] %s", errx.Error()))
		return types.User{}, errx
	}

	return types.User{
		ID:            userID,
		Email:         emailID,
		KeyHash:       hashStr,
		EncryptionKey: encryptionKeyStr,
	}, nil
}
