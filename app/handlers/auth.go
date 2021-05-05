package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RefreshTokenHandler(jwtCfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("refresh_token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var claims types.AccessTokenClaims
		var err error
		if claims, err = utils.ValidateJWT(token, jwtCfg.GetSecret(), lgr); err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		accessTokenClaims := types.AccessTokenClaims{
			UserID:        claims.UserID,
			EncryptionKey: claims.EncryptionKey,
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}

		tkn, err := utils.IssueJWT(accessTokenClaims, jwtCfg.GetSecret(), lgr)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: tkn,
			RefreshToken:        token,
		}

		d, err := json.Marshal(resp)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(d)
		if err != nil {
			// TODO: Add Logging
			return
		}
	}
}
