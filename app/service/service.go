package service

import (
	"encoding/base64"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type Service interface {
	CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, error)
	GetUser(emailID string) (types.User, error)
}

type dbServiceInstance struct {
	db  *database.DB
	lgr *zap.Logger
}

func NewDBService(db *database.DB, lgr *zap.Logger) Service {
	return &dbServiceInstance{
		db:  db,
		lgr: lgr,
	}
}

func (s *dbServiceInstance) GetUser(emailID string) (types.User, error) {
	usr, err := s.db.Users.Get(emailID)
	if err != nil {
		// TODO: Add Logging
		// TODO: Check for no records error
		return types.User{}, err
	}

	return usr, nil
}

func (s *dbServiceInstance) CreateUser(emailID string, encryptionKey []byte, keyHash [32]byte) (types.User, error) {
	encryptionKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)
	hashStr := base64.StdEncoding.EncodeToString(keyHash[:])

	userID, err := s.db.Users.Create(emailID, encryptionKeyStr, hashStr)
	if err != nil {
		s.lgr.Sugar().Error(err)
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
