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

type NotesService interface {
	GetAll(claims types.AccessTokenClaims) ([]types.Note, *erx.Erx)
	Create(name string, data string, folderID types.FolderID, claims types.AccessTokenClaims) (types.NoteID, *erx.Erx)
	Update(name string, data string, folderID types.FolderID, noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx
	Delete(noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx
}

type notes struct {
	db  *database.DB
	lgr *zap.Logger
}

func (n *notes) GetAll(claims types.AccessTokenClaims) ([]types.Note, *erx.Erx) {
	notesList, errx := n.db.Notes.GetAll(claims.UserID)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [GetAll] [GetAll] %s", errx.String()))
		return nil, errx
	}

	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [GetAll] [NewCipher] %s", err.Error()))
		return nil, erx.WithArgs(err, erx.SeverityDebug)
	}

	for ind, note := range notesList {
		name, err := base64.StdEncoding.DecodeString(note.Name)
		if err != nil {
			n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [GetAll] [DecodeString] [Name] %s", err.Error()))
			return nil, erx.WithArgs(err, erx.SeverityDebug)
		}
		note.Name = string(utils.CFBDecrypt(name, blockCipher))

		data, err := base64.StdEncoding.DecodeString(note.Data)
		if err != nil {
			n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [GetAll] [DecodeString] [Data] %s", err.Error()))
			return nil, erx.WithArgs(err, erx.SeverityDebug)
		}
		note.Data = string(utils.CFBDecrypt(data, blockCipher))

		notesList[ind] = note
	}

	return notesList, nil
}

func (n *notes) Create(name string, data string, folderID types.FolderID, claims types.AccessTokenClaims) (types.NoteID, *erx.Erx) {
	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Create] [NewCipher] %s", err.Error()))
		return 0, erx.WithArgs(err, erx.SeverityDebug)
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Create] [CFBEncrypt] [Name] %s", err.Error()))
		return 0, erx.WithArgs(err, erx.SeverityDebug)
	}

	encryptedData, err := utils.CFBEncrypt([]byte(data), blockCipher)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Create] [CFBEncrypt] [Data] %s", err.Error()))
		return 0, erx.WithArgs(err, erx.SeverityDebug)
	}

	noteID, errx := n.db.Notes.Create(base64.StdEncoding.EncodeToString(encryptedName),
		base64.StdEncoding.EncodeToString(encryptedData), folderID, claims.UserID)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Create] [Create] %s", errx.String()))
		return 0, errx
	}

	return noteID, nil
}

func (n *notes) Update(name string, data string, folderID types.FolderID, noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx {
	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Update] [NewCipher] %s", err.Error()))
		return erx.WithArgs(err, erx.SeverityDebug)
	}

	encryptedName, err := utils.CFBEncrypt([]byte(name), blockCipher)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Update] [CFBEncrypt] [Name] %s", err.Error()))
		return erx.WithArgs(err, erx.SeverityDebug)
	}

	encryptedData, err := utils.CFBEncrypt([]byte(data), blockCipher)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Update] [CFBEncrypt] [Data] %s", err.Error()))
		return erx.WithArgs(err, erx.SeverityDebug)
	}

	errx := n.db.Notes.Update(types.Note{
		ID:       noteID,
		FolderID: folderID,
		Data:     base64.StdEncoding.EncodeToString(encryptedData),
		Name:     base64.StdEncoding.EncodeToString(encryptedName),
	}, claims.UserID)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Update] [Uodate] %s", errx.String()))
		return errx
	}

	return nil
}

func (n *notes) Delete(noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx {
	errx := n.db.Notes.Delete(noteID, claims.UserID)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Delete] [Delete] %s", errx.String()))
		return errx
	}
	return nil
}
