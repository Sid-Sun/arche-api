package database

import "github.com/sid-sun/arche-api/app/types"

type FoldersTable interface {
	Get(userID types.UserID) ([]types.Folder, error)
	Delete(folderID types.FolderID, UserID types.UserID) error
	Create(name string, userID types.UserID) (types.FolderID, error)
}
