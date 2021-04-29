package types

import "github.com/dgrijalva/jwt-go"

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	AuthenticationToken string `json:"authentication_token"`
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

type DeleteFolderRequest struct {
	FolderID FolderID `json:"folder_id"`
}
