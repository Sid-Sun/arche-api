package types

import "github.com/dgrijalva/jwt-go"

type CreateUserRequest struct {
	Email                   string `json:"email"`
	Password                string `json:"password"`
	VerificationCallbackURL string `json:"verification_callback_url"`
}

type CreateUserResponse struct {
	UserCreated           bool `json:"user_created"`
	VerificationEmailSent bool `json:"verification_email_sent"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AuthenticationToken string `json:"authentication_token,omitempty"`
	RefreshToken        string `json:"refresh_token,omitempty"`
	VerificationPending bool   `json:"verification_pending"`
}

type ActivateUserRequest struct {
	VerificationToken string `json:"verification_token"`
}

type ActivateUserResponse struct {
	VerificationPending bool `json:"verification_pending"`
}

type ResendVerificationRequest struct {
	Email                   string `json:"email"`
	VerificationCallbackURL string `json:"verification_callback_url"`
}

type ResendVerificationResponse struct {
	VerificationEmailSent bool `json:"verification_email_sent"`
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
