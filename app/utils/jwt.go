package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"time"
)

func IssueJWT(claims types.AccessTokenClaims, secret string, lgr *zap.Logger) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tkn.SignedString([]byte(secret))
	if err != nil {
		// TODO: Add logging
		return "", err
	}
	return token, nil
}

func ValidateJWT(tkn string, secret string, lgr *zap.Logger) (types.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tkn, &types.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		// TODO: Add logging
		return types.AccessTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*types.AccessTokenClaims); ok && token.Valid {
		return *claims, nil
	}

	return types.AccessTokenClaims{}, errors.New("token not okay or invalid or incorrect token claim")
}

func IssueTokens(userID types.UserID, key []byte, cfg *config.JWTConfig, lgr *zap.Logger) (accessToken string, refreshToken string, err error) {
	refreshClaims := types.AccessTokenClaims{
		UserID:        userID,
		EncryptionKey: key,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	accessClaims := types.AccessTokenClaims{
		UserID:        userID,
		EncryptionKey: key,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	refreshToken, err = IssueJWT(refreshClaims, cfg.GetSecret(), lgr)
	if err != nil {
		// TODO: Add Logging
		return "", "", err
	}

	accessToken, err = IssueJWT(accessClaims, cfg.GetSecret(), lgr)
	if err != nil {
		// TODO: Add Logging
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
