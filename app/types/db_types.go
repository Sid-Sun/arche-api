package types

type UserID int
type FolderID int
type NoteID int

type User struct {
	ID            UserID
	Email         string
	EncryptionKey string
	KeyHash       string
}

type Folder struct {
	ID     FolderID
	UserID UserID
	Name   string
}

type Note struct {
	ID       NoteID
	FolderID FolderID
	Data     string
	Name     string
}
