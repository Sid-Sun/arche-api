package database

import (
	"database/sql"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type FoldersTable interface {
	Get(userID types.UserID) ([]types.Folder, error)
	Delete(folderID types.FolderID, UserID types.UserID) error
	Create(name string, userID types.UserID) (types.FolderID, error)
}

type folders struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (f *folders) Get(userID types.UserID) ([]types.Folder, error) {
	query := `SELECT folder_id, name FROM folders WHERE user_id=@user_id`

	rows, err := f.db.Query(query, sql.Named("user_id", userID))
	if err != nil {
		// TODO: Add Logging
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// TODO: Add Logging
		}
	}(rows)
	folders := *new([]types.Folder)

	for rows.Next() {
		var folderID types.FolderID
		var name string

		err = rows.Scan(&folderID, &name)
		if err != nil {
			// TODO: Add Logging
			return nil, err
		}

		folders = append(folders, types.Folder{
			ID:     folderID,
			UserID: userID,
			Name:   name,
		})
	}

	if err := rows.Err(); err != nil {
		// TODO: Add Logging
		return nil, err
	}

	return folders, nil
}
func (f *folders) Delete(folderID types.FolderID, userID types.UserID) error {
	query := `DELETE FROM folders WHERE folder_id=@folder_id AND user_id=@user_id`

	row := f.db.QueryRow(query, sql.Named("@folder_id", folderID), sql.Named("user_id", userID))
	err := row.Err()
	if err != nil {
		// TODO: Add Logging
		return err
	}

	err = row.Scan(&folderID)
	if err != nil {
		// TODO: Add Logging
		return err
	}

	return nil
}

func (f *folders) Create(name string, userID types.UserID) (types.FolderID, error) {
	query := `INSERT INTO folders (user_id, name) OUTPUT inserted.folder_id VALUES (@user_id, @name)`

	row := f.db.QueryRow(query, sql.Named("user_id", userID), sql.Named("name", name))
	err := row.Err()
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	var folderID types.FolderID
	err = row.Scan(&folderID)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	return folderID, nil
}
