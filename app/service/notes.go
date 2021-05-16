package service

import (
	"crypto/aes"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type NotesService interface {
	Get(noteID types.NoteID, claims types.AccessTokenClaims) (types.Note, *erx.Erx)
	GetAll(claims types.AccessTokenClaims) ([]types.Note, *erx.Erx)
	Create(name string, data string, folderID types.FolderID, claims types.AccessTokenClaims) (types.NoteID, *erx.Erx)
	Update(name string, data string, folderID types.FolderID, noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx
	Delete(noteID types.NoteID, claims types.AccessTokenClaims) *erx.Erx
}

type notes struct {
	db  *database.DB
	lgr *zap.Logger
}

func (n *notes) Get(noteID types.NoteID, claims types.AccessTokenClaims) (types.Note, *erx.Erx) {
	note, errx := n.db.Notes.Get(noteID, claims.UserID)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Get] [Get] %s", errx.String()))
		return types.Note{}, errx
	}

	blockCipher, err := aes.NewCipher(claims.EncryptionKey)
	if err != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Get] [NewCipher] %s", err.Error()))
		return types.Note{}, erx.WithArgs(err, erx.SeverityDebug)
	}

	note, errx = decryptNote(note, blockCipher, n.lgr)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Get] [decryptNote] %s", errx.String()))
		return types.Note{}, errx
	}

	return note, nil
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
		note, errx := decryptNote(note, blockCipher, n.lgr)
		if errx != nil {
			n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [GetAll] [decryptNote] %s", errx.String()))
			return nil, errx
		}
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

	// Create note with zero ID as encryptNote needs note
	// It does not operate on ID
	note := types.Note{
		FolderID: folderID,
		ID:       0,
		Name:     name,
		Data:     data,
	}
	note, errx := encryptNote(note, blockCipher, n.lgr)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Create] [encryptNote] %s", errx.String()))
		return 0, errx
	}

	noteID, errx := n.db.Notes.Create(note.Name, note.Data, folderID, claims.UserID)
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

	// Create note with zero ID as encryptNote needs note
	// It does not operate on ID
	note := types.Note{
		FolderID: folderID,
		ID:       noteID,
		Name:     name,
		Data:     data,
	}
	note, errx := encryptNote(note, blockCipher, n.lgr)
	if errx != nil {
		n.lgr.Debug(fmt.Sprintf("[Service] [Notes] [Update] [encryptNote] %s", errx.String()))
		return errx
	}

	errx = n.db.Notes.Update(note, claims.UserID)
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
