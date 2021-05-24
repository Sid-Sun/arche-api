package types

import "github.com/dgrijalva/jwt-go"

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	AuthenticationToken string `json:"authentication_token"`
	RefreshToken        string `json:"refresh_token"`
}

type AccessTokenClaims struct {
	UserID        UserID `json:"user_id"`
	EncryptionKey []byte `json:"encryption_key"`
	jwt.StandardClaims
}

type CreateFolderRequest struct {
	Name string `json:"name"`
}

type CreateFolderResponse struct {
	Name     string   `json:"name"`
	FolderID FolderID `json:"folder_id"`
}

type DeleteFolderResponse DeleteFolderRequest
type DeleteFolderRequest struct {
	FolderID FolderID `json:"folder_id"`
}

type CreateNoteResponse struct {
	NoteID NoteID `json:"note_id"`
}

type DeleteNoteResponse DeleteNoteRequest
type DeleteNoteRequest struct {
	NoteID NoteID `json:"note_id"`
}

type CreateNoteRequest struct {
	Name     string   `json:"name"`
	Data     string   `json:"data"`
	FolderID FolderID `json:"folder_id"`
}

type UpdateNoteResponse UpdateNoteRequest
type UpdateNoteRequest struct {
	Name     string   `json:"name"`
	Data     string   `json:"data"`
	NoteID   NoteID   `json:"note_id"`
	FolderID FolderID `json:"folder_id"`
}
