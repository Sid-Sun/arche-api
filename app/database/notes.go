package database

import "github.com/sid-sun/arche-api/app/types"

type NotesTable interface {
	GetAll(userID types.UserID) ([]types.Note, error)
	Create(name string, data string, folderID types.FolderID, userID types.UserID) (types.NoteID, error)
	Update(name string, data string, folderID types.FolderID, userID types.UserID) error
	Delete(noteID types.NoteID, userID types.UserID) error
}
