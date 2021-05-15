package service

import (
	"crypto/aes"
	"encoding/base64"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

type NotesService interface {
	GetAll(claims types.AccessTokenClaims) ([]types.Note, error)
	Create(name string, data string, folderID types.FolderID, claims types.AccessTokenClaims) (types.NoteID, error)
	Update(name string, data string, folderID types.FolderID, noteID types.NoteID, claims types.AccessTokenClaims) error
	Delete(noteID types.NoteID, claims types.AccessTokenClaims) error
}

type notes struct {
	db  *database.DB
	lgr *zap.Logger
}

func (n *notes) GetAll(claims types.AccessTokenClaims) ([]types.Note, error) {
	notesList, errx := n.db.Notes.GetAll(claims.UserID)
	if errx != nil {
		return nil, errx
	}

	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		// TODO: Add Logging
		return nil, err
	}

	for ind, note := range notesList {
		name, err := base64.StdEncoding.DecodeString(note.Name)
		if err != nil {
			// TODO: Add Logging
			return nil, err
		}
		note.Name = string(utils.CFBDecrypt(name, blockCipher))

		data, err := base64.StdEncoding.DecodeString(note.Data)
		if err != nil {
			// TODO: Add Logging
			return nil, err
		}
		note.Data = string(utils.CFBDecrypt(data, blockCipher))

		notesList[ind] = note
	}

	return notesList, err
}

func (n *notes) Create(name string, data string, folderID types.FolderID, claims types.AccessTokenClaims) (types.NoteID, error) {
	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	encryptedData, err := utils.CFBEncrypt([]byte(data), blockCipher)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	noteID, err := n.db.Notes.Create(base64.StdEncoding.EncodeToString(encryptedName),
		base64.StdEncoding.EncodeToString(encryptedData), folderID, claims.UserID)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	return noteID, nil
}

func (n *notes) Update(name string, data string, folderID types.FolderID, noteID types.NoteID, claims types.AccessTokenClaims) error {
	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		// TODO: Add Logging
		return err
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		// TODO: Add Logging
		return err
	}

	encryptedData, err := utils.CFBEncrypt([]byte(data), blockCipher)
	if err != nil {
		// TODO: Add Logging
		return err
	}

	err = n.db.Notes.Update(types.Note{
		ID:       noteID,
		FolderID: folderID,
		Data:     base64.StdEncoding.EncodeToString(encryptedData),
		Name:     base64.StdEncoding.EncodeToString(encryptedName),
	}, claims.UserID)
	if err != nil {
		// TODO: Add Logging
		return err
	}

	return nil
}

func (n *notes) Delete(noteID types.NoteID, claims types.AccessTokenClaims) error {
	err := n.db.Notes.Delete(noteID, claims.UserID)
	if err != nil {
		// TODO: Add Logging
		return err
	}
	return nil
}
