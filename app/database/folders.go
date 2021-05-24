package database

import (
	"database/sql"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type FoldersTable interface {
	Get(folderID types.FolderID, userID types.UserID) ([]types.FolderContent, *erx.Erx)
	GetAll(userID types.UserID) ([]types.Folder, *erx.Erx)
	Delete(folderID types.FolderID, UserID types.UserID) *erx.Erx
	Create(name string, userID types.UserID) (types.FolderID, *erx.Erx)
}

type folders struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (f *folders) Get(folderID types.FolderID, userID types.UserID) ([]types.FolderContent, *erx.Erx) {
	query := `SELECT note_id, name FROM notes WHERE folder_id=(SELECT folder_id FROM folders WHERE user_id=@user_id AND folder_id=@folder_id)`

	rows, err := f.db.Query(query, sql.Named("user_id", userID), sql.Named("folder_id", folderID))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Get] [Query] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Get] [Query] %d : %s", sqlErr.Number, sqlErr.Error()))
		return nil, errx
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Get] [Close] %s", err.Error()))
		}
	}(rows)
	folderContents := *new([]types.FolderContent)

	for rows.Next() {
		var noteID types.NoteID
		var name string

		err = rows.Scan(&noteID, &name)
		if err != nil {
			sqlErr, errx := checkForSQLError(err)
			if sqlErr != nil {
				errx = erx.WithArgs(errx, erx.SeverityError)
				f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Get] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
				return nil, errx
			}
			f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Get] [Scan] %s", err.Error()))
			return nil, errx
		}

		folderContents = append(folderContents, types.FolderContent{
			NoteID: noteID,
			Name: name,
		})
	}

	if err := rows.Err(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Get] [Err] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Get] [Err] %s", err.Error()))
		return nil, errx
	}

	return folderContents, nil
}

func (f *folders) GetAll(userID types.UserID) ([]types.Folder, *erx.Erx) {
	query := `SELECT folder_id, name FROM folders WHERE user_id=@user_id`

	rows, err := f.db.Query(query, sql.Named("user_id", userID))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [GetAll] [Query] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [GetAll] [Query] %d : %s", sqlErr.Number, sqlErr.Error()))
		return nil, errx
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [GetAll] [Close] %s", err.Error()))
		}
	}(rows)
	folders := *new([]types.Folder)

	for rows.Next() {
		var folderID types.FolderID
		var name string

		err = rows.Scan(&folderID, &name)
		if err != nil {
			sqlErr, errx := checkForSQLError(err)
			if sqlErr != nil {
				errx = erx.WithArgs(errx, erx.SeverityError)
				f.lgr.Error(fmt.Sprintf("[Database] [Folders] [GetAll] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
				return nil, errx
			}
			f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [GetAll] [Scan] %s", err.Error()))
			return nil, errx
		}

		folders = append(folders, types.Folder{
			FolderID: folderID,
			UserID:   userID,
			Name:     name,
		})
	}

	if err := rows.Err(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [GetAll] [Err] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [GetAll] [Err] %s", err.Error()))
		return nil, errx
	}

	return folders, nil
}

func (f *folders) Delete(folderID types.FolderID, userID types.UserID) *erx.Erx {
	query := `DELETE FROM folders WHERE folder_id=@folder_id AND user_id=@user_id`

	res, err := f.db.Exec(query, sql.Named("folder_id", folderID), sql.Named("user_id", userID))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Delete] [Exec] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Delete] [Exec] %s", err.Error()))
		return errx
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Delete] [RowsAffected] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Delete] [RowsAffected] %s", err.Error()))
		return errx
	}

	if count == 0 {
		return erx.WithArgs(custom_errors.NoRowsAffected, erx.SeverityInfo)
	}

	return nil
}

func (f *folders) Create(name string, userID types.UserID) (types.FolderID, *erx.Erx) {
	query := `INSERT INTO folders (user_id, name) OUTPUT inserted.folder_id VALUES (@user_id, @name)`

	row := f.db.QueryRow(query, sql.Named("user_id", userID), sql.Named("name", name))
	err := row.Err()
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Create] [QueryRow] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Create] [QueryRow] %s", err.Error()))
		return 0, errx
	}

	var folderID types.FolderID
	err = row.Scan(&folderID)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			f.lgr.Error(fmt.Sprintf("[Database] [Folders] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		f.lgr.Debug(fmt.Sprintf("[Database] [Folders] [Create] [Scan] %s", err.Error()))
		return 0, errx
	}

	return folderID, nil
}
