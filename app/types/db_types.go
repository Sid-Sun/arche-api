package types

type UserID int
type FolderID int
type NoteID int

type User struct {
	ID              UserID `json:"user_id"`
	Email           string `json:"email"`
	KeyHash         string `json:"key_hash"`
	EncryptionKey   string `json:"encryption_key"`
	VerificationKey string `json:"verification_key"`
	Verified        bool   `json:"verified"`
}

type Folder struct {
	FolderID FolderID `json:"folder_id"`
	UserID   UserID   `json:"user_id"`
	Name     string   `json:"name"`
}

type Note struct {
	NoteID   NoteID   `json:"note_id"`
	FolderID FolderID `json:"folder_id"`
	Data     string   `json:"data"`
	Name     string   `json:"name"`
}
