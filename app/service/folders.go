package service

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

type FoldersService interface {
	Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, *erx.Erx)
	GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, *erx.Erx)
	Delete(folderID types.FolderID, userClaims types.AccessTokenClaims) *erx.Erx
}

type folders struct {
	db  *database.DB
	lgr *zap.Logger
}

func (f folders) Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, *erx.Erx) {
	blockCipher, err := aes.NewCipher(userClaims.EncryptionKey)
	if err != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Create] [NewCipher] %s", err.Error()))
		return 0, erx.WithArgs(err, erx.SeverityDebug)
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Create] [CFBEncrypt] [Name] %s", err.Error()))
		return 0, erx.WithArgs(err, erx.SeverityDebug)
	}

	folderID, errx := f.db.Folders.Create(base64.StdEncoding.EncodeToString(encryptedName), userClaims.UserID)
	if errx != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Create] [Create] %s", errx.String()))
		// TODO: Add Logging
		return 0, errx
	}

	return folderID, nil
}

func (f folders) GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, *erx.Erx) {
	fldrs, errx := f.db.Folders.Get(userClaims.UserID)
	if errx != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [GetAll] [Get] %s", errx.String()))
		return []types.Folder{}, errx
	}

	blockCipher, err := aes.NewCipher(userClaims.EncryptionKey)
	if err != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [GetAll] [NewCipher] %s", err.Error()))
		return []types.Folder{}, erx.WithArgs(err, erx.SeverityDebug)
	}

	for index, folder := range fldrs {
		name, err := base64.StdEncoding.DecodeString(folder.Name)
		if err != nil {
			f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [GetAll] [DecodeString] %s", err.Error()))
			return []types.Folder{}, erx.WithArgs(err, erx.SeverityDebug)
		}

		folder.Name = string(utils.CFBDecrypt(name, blockCipher))
		fldrs[index] = folder
	}

	return fldrs, nil
}

func (f folders) Delete(folderID types.FolderID, userClaims types.AccessTokenClaims) *erx.Erx {
	errx := f.db.Folders.Delete(folderID, userClaims.UserID)
	if errx != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Delete] [Delete] %s", errx.String()))
		return errx
	}
	return nil
}
