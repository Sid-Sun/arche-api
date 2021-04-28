package types

import "github.com/dgrijalva/jwt-go"

type UserRequest struct {
	Email    string
	Password string
}

type AccessTokenClaims struct {
	UserID        UserID `json:"user_id"`
	EncryptionKey []byte `json:"encryption_key"`
	jwt.StandardClaims
}
