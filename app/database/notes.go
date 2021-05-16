package database

import (
	"database/sql"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type NotesTable interface {
	Get(noteID types.NoteID, userID types.UserID) (types.Note, *erx.Erx)
	GetAll(userID types.UserID) ([]types.Note, *erx.Erx)
	Create(name string, data string, folderID types.FolderID, userID types.UserID) (types.NoteID, *erx.Erx)
	Update(note types.Note, userID types.UserID) *erx.Erx
	Delete(noteID types.NoteID, userID types.UserID) *erx.Erx
}

type notes struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (n *notes) Get(noteID types.NoteID, userID types.UserID) (types.Note, *erx.Erx) {
	query := `SELECT notes.name, notes.data, f.folder_id FROM notes INNER JOIN folders AS f ON (f.folder_id = notes.folder_id) WHERE user_id=@user_id AND note_id=@note_id`

	row := n.db.QueryRow(query, sql.Named("user_id", userID), sql.Named("note_id", noteID))
	err := row.Err()
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Get] [QueryRow] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return types.Note{}, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Get] [QueryRow] %s", err.Error()))
		return types.Note{}, errx
	}

	var folderID types.FolderID
	var name, data string
	err = row.Scan(&name, &data, &folderID)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Get] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return types.Note{}, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Get] [Scan] %s", err.Error()))
		return types.Note{}, errx
	}

	return types.Note{
		ID: noteID,
		Name: name,
		Data: data,
		FolderID: folderID,
	}, nil
}

func (n *notes) GetAll(userID types.UserID) ([]types.Note, *erx.Erx) {
	query := `SELECT notes.note_id, notes.name, notes.data, notes.folder_id
FROM notes INNER JOIN folders AS f ON (f.folder_id = notes.folder_id) WHERE user_id=@userID`

	rows, err := n.db.Query(query, sql.Named("userID", userID))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [GetAll] [Query] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [GetAll] [Query] %s", err.Error()))
		return nil, errx
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [GetAll] [Close] %s", err.Error()))
		}
	}(rows)
	notesSlice := *new([]types.Note)

	for rows.Next() {
		var noteID types.NoteID
		var folderID types.FolderID
		var data string
		var name string

		err = rows.Scan(&noteID, &name, &data, &folderID)
		if err != nil {
			sqlErr, errx := checkForSQLError(err)
			if sqlErr != nil {
				errx = erx.WithArgs(errx, erx.SeverityError)
				n.lgr.Error(fmt.Sprintf("[Database] [Notes] [GetAll] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
				return nil, errx
			}
			n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [GetAll] [Scan] %s", err.Error()))
			return nil, errx
		}

		notesSlice = append(notesSlice, types.Note{
			ID:       noteID,
			FolderID: folderID,
			Data:     data,
			Name:     name,
		})

	}

	if err := rows.Err(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [GetAll] [Err] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return nil, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [GetAll] [Err] %s", err.Error()))
		return nil, errx
	}

	return notesSlice, nil
}

func (n *notes) Create(name string, data string, folderID types.FolderID, userID types.UserID) (types.NoteID, *erx.Erx) {
	query := `INSERT INTO notes (data, name, folder_id) OUTPUT inserted.note_id 
VALUES (@data, @name, (SELECT folder_id FROM folders WHERE user_id=@userID AND  folder_id=@folderID))`

	row := n.db.QueryRow(query, sql.Named("data", data), sql.Named("name", name),
		sql.Named("userID", userID), sql.Named("folderID", folderID))
	err := row.Err()
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Create] [QueryRow] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Create] [QueryRow] %s", err.Error()))
		return 0, errx
	}

	var noteID types.NoteID
	err = row.Scan(&noteID)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Create] [Scan] %s", err.Error()))
		return 0, errx
	}

	return noteID, nil
}

func (n *notes) Update(note types.Note, userID types.UserID) *erx.Erx {
	query := `UPDATE notes SET name = @name, data = @data 
WHERE note_id = @noteID AND folder_id = (SELECT folder_id FROM folders WHERE folder_id = (SELECT notes.folder_id FROM notes WHERE note_id = @noteID) AND user_id = @userID)`

	res, err := n.db.Exec(query, sql.Named("name", note.Name), sql.Named("data", note.Data),
		sql.Named("noteID", note.ID), sql.Named("userID", userID))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Update] [Exec] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Update] [Exec] %s", err.Error()))
		return errx
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Update] [RowsAffected] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Update] [RowsAffected] %s", err.Error()))
		return errx
	}

	if count == 0 {
		return erx.WithArgs(custom_errors.NoRowsAffected, erx.SeverityInfo)
	}

	return nil
}

func (n *notes) Delete(noteID types.NoteID, userID types.UserID) *erx.Erx {
	query := `DELETE FROM notes WHERE note_id = @noteID AND folder_id = 
                                              (SELECT folder_id FROM folders WHERE folder_id = (
                                                  SELECT notes.folder_id FROM notes WHERE note_id = @noteID) AND user_id = @userID)`

	res, err := n.db.Exec(query, sql.Named("noteID", noteID), sql.Named("userID", userID))

	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Delete] [Exec] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Delete] [Exec] %s", err.Error()))
		return errx
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			n.lgr.Error(fmt.Sprintf("[Database] [Notes] [Delete] [RowsAffected] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		n.lgr.Debug(fmt.Sprintf("[Database] [Notes] [Delete] [RowsAffected] %s", err.Error()))
		return errx
	}

	if count == 0 {
		return erx.WithArgs(custom_errors.NoRowsAffected, erx.SeverityInfo)
	}

	return nil
}
