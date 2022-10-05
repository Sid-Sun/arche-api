package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/initializers"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

type UsersService interface {
	SendVerificationEmail(emailID string, verificationString string, callbackURL string, veCfg *config.VerificationEmailConfig) *erx.Erx
	CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte, vetkn string) (types.User, *erx.Erx)
	ActivateUser(verificationString string) *erx.Erx
	GetVerificationStatus(emailID string) (bool, string, *erx.Erx)
	UpdateVerificationToken(email string, token string) *erx.Erx
	GetUser(emailID string) (types.User, *erx.Erx)
}

type users struct {
	db         *database.DB
	lgr        *zap.Logger
	mailClient initializers.MailClient
}

func (u *users) GetUser(emailID string) (types.User, *erx.Erx) {
	usr, errx := u.db.Users.Get(emailID)
	if errx != nil {
		(*u).lgr.Debug(fmt.Sprintf("[Service] [Users] [GetUser] [Get] %s", errx.Error()))
		return types.User{}, errx
	}

	return usr, nil
}

func (u *users) CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte, vetkn string) (types.User, *erx.Erx) {
	encryptionKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)
	hashStr := base64.StdEncoding.EncodeToString(keyHash[:])

	userID, errx := u.db.Users.Create(emailID, encryptionKeyStr, hashStr, vetkn)
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

func (u *users) ActivateUser(verificationString string) *erx.Erx {
	errx := u.db.Users.VerifyUser(verificationString)
	if errx != nil {
		(*u).lgr.Debug(fmt.Sprintf("[Service] [Users] [GetUser] [Get] %s", errx.Error()))
		return errx
	}
	return nil
}

func (u *users) GetVerificationStatus(emailID string) (bool, string, *erx.Erx) {
	return u.db.Users.GetVerificationStatus(emailID)
}

func (u *users) SendVerificationEmail(emailID string, verificationString string, callbackURL string, veCfg *config.VerificationEmailConfig) *erx.Erx {
	callbackURL = fmt.Sprintf("%s?verificationToken=%s", callbackURL, verificationString)

	msg := u.mailClient.NewMessage(veCfg.GetSenderEmail(u.mailClient.Domain()), veCfg.GetSubject(), veCfg.GetBody(callbackURL), emailID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := u.mailClient.Send(ctx, msg)
	if err != nil {
		return erx.WithArgs(err, erx.SeverityDebug)
	}

	return nil
}

func (u *users) UpdateVerificationToken(email string, token string) *erx.Erx {
	return u.db.Users.UpdateVerificationToken(email, token)
}
