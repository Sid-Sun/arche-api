package database

import (
	"database/sql"
	"errors"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type NotesTable interface {
	GetAll(userID types.UserID) ([]types.Note, error)
	Create(name string, data string, folderID types.FolderID, userID types.UserID) (types.NoteID, error)
	Update(note types.Note, userID types.UserID) error
	Delete(noteID types.NoteID, userID types.UserID) error
}

type notes struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (n *notes) GetAll(userID types.UserID) ([]types.Note, error) {
	query := `SELECT notes.note_id, notes.name, notes.data, notes.folder_id
FROM notes INNER JOIN folders AS f ON (f.folder_id = notes.folder_id) WHERE user_id=@userID`

	rows, err := n.db.Query(query, sql.Named("userID", userID))
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
	notesSlice := *new([]types.Note)

	for rows.Next() {
		var noteID types.NoteID
		var folderID types.FolderID
		var data string
		var name string

		err = rows.Scan(&noteID, &name, &data, &folderID)
		if err != nil {
			// TODO: Add Logging
			return nil, err
		}

		notesSlice = append(notesSlice, types.Note{
			ID:       noteID,
			FolderID: folderID,
			Data:     data,
			Name:     name,
		})

	}

	if err := rows.Err(); err != nil {
		// TODO: Add Logging
		return nil, err
	}

	return notesSlice, nil
}

func (n *notes) Create(name string, data string, folderID types.FolderID, userID types.UserID) (types.NoteID, error) {
	query := `INSERT INTO notes (data, name, folder_id) OUTPUT inserted.note_id 
VALUES (@data, @name, (SELECT folder_id FROM folders WHERE user_id=@userID AND  folder_id=@folderID))`

	row := n.db.QueryRow(query, sql.Named("data", data), sql.Named("name", name),
		sql.Named("userID", userID), sql.Named("folderID", folderID))
	err := row.Err()
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	var noteID types.NoteID
	err = row.Scan(&noteID)
	if err != nil {
		// TODO: Add Logging
		return 0, err
	}

	return noteID, nil
}

func (n *notes) Update(note types.Note, userID types.UserID) error {
	query := `UPDATE notes SET name = @name, data = @data 
WHERE note_id = @noteID AND folder_id = (SELECT folder_id FROM folders WHERE folder_id = (SELECT notes.folder_id FROM notes WHERE note_id = @noteID) AND user_id = @userID)`

	res, err := n.db.Exec(query, sql.Named("name", note.Name), sql.Named("data", note.Data),
		sql.Named("noteID", note.ID), sql.Named("userID", userID))

	if err != nil {
		// TODO: Add Logging
		return err
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		if err != nil {
			// TODO: Add Logging
			// TODO: Add Error Handling
			return err
		}
		return errors.New("no records were modified")
	}

	return nil
}

func (n *notes) Delete(noteID types.NoteID, userID types.UserID) error {
	query := `DELETE FROM notes WHERE note_id = @noteID AND folder_id = 
                                              (SELECT folder_id FROM folders WHERE folder_id = (
                                                  SELECT notes.folder_id FROM notes WHERE note_id = @noteID) AND user_id = @userID)`

	res, err := n.db.Exec(query, sql.Named("noteID", noteID), sql.Named("userID", userID))

	if err != nil {
		// TODO: Add Logging
		return err
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		if err != nil {
			// TODO: Add Logging
			// TODO: Add Error Handling
			return err
		}
		return errors.New("no records were deleted")
	}

	return nil
}
