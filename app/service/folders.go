package service

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

type FoldersService interface {
	Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, *erx.Erx)
	Get(folderID types.FolderID, userClaims types.AccessTokenClaims) ([]types.FolderContent, *erx.Erx)
	GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, *erx.Erx)
	Delete(folderID types.FolderID, userClaims types.AccessTokenClaims) *erx.Erx
}

type folders struct {
	db  *database.DB
	lgr *zap.Logger
}

func (f folders) Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, *erx.Erx) {
	fldrs, errx := f.GetAll(userClaims)
	if errx != nil {
		return 0, errx
	}
	for _, flder := range fldrs {
		if name == flder.Name {
			return 0, erx.WithArgs(errors.New("folder already exists"), custom_errors.DuplicateRecordInsertion, erx.SeverityInfo)
		}
	}

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
		return 0, errx
	}

	return folderID, nil
}

func (f folders) Get(folderID types.FolderID, userClaims types.AccessTokenClaims) ([]types.FolderContent, *erx.Erx) {
	contents, errx := f.db.Folders.Get(folderID, userClaims.UserID)
	if errx != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Get] [Get] %s", errx.String()))
		return nil, errx
	}

	blockCipher, err := aes.NewCipher(userClaims.EncryptionKey)
	if err != nil {
		f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Get] [NewCipher] %s", err.Error()))
		return nil, erx.WithArgs(err, erx.SeverityDebug)
	}

	for index, content := range contents {
		name, err := base64.StdEncoding.DecodeString(content.Name)
		if err != nil {
			f.lgr.Debug(fmt.Sprintf("[Service] [Folders] [Get] [DecodeString] %s", err.Error()))
			return nil, erx.WithArgs(err, erx.SeverityDebug)
		}

		content.Name = string(utils.CFBDecrypt(name, blockCipher))
		contents[index] = content
	}

	return contents, nil
}

func (f folders) GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, *erx.Erx) {
	fldrs, errx := f.db.Folders.GetAll(userClaims.UserID)
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
