package service

import (
	"crypto/aes"
	"encoding/base64"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

type FoldersService interface {
	Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, error)
	GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, error)
	Delete(folderID types.FolderID, userClaims types.AccessTokenClaims) error
}

type folders struct {
	db  *database.DB
	lgr *zap.Logger
}

func (f folders) Create(name string, userClaims types.AccessTokenClaims) (types.FolderID, error) {
	blockCipher, err := aes.NewCipher(userClaims.EncryptionKey)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	folderID, err := f.db.Folders.Create(base64.StdEncoding.EncodeToString(encryptedName), userClaims.UserID)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	return folderID, nil
}

func (f folders) GetAll(userClaims types.AccessTokenClaims) ([]types.Folder, error) {
	fldrs, err := f.db.Folders.Get(userClaims.UserID)
	if err != nil {
		// TODO: Add Logging
		return []types.Folder{}, err
	}

	blockCipher, err := aes.NewCipher(userClaims.EncryptionKey)
	if err != nil {
		// TODO: Add Logging
		return []types.Folder{}, err
	}

	for index, folder := range fldrs {
		name, err := base64.StdEncoding.DecodeString(folder.Name)
		if err != nil {
			// TODO: Add Logging
			return []types.Folder{}, err
		}

		folder.Name = string(utils.CFBDecrypt(name, blockCipher))
		fldrs[index] = folder
	}

	return fldrs, nil
}

func (f folders) Delete(folderID types.FolderID, userClaims types.AccessTokenClaims) error {
	err := f.db.Folders.Delete(folderID, userClaims.UserID)
	if err != nil {
		// TODO: Add Logging
		return err
	}
	return nil
}
