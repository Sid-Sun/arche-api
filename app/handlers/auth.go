package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sid-sun/arche-api/app/http/resperr"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
)

func RefreshTokenHandler(jwtCfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("refresh_token")
		if token == "" {
			lgr.Debug("[Handlers] [RefreshTokenHandler] [Get] token not in request headers")
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "refresh token was not sent"), w, lgr)
			return
		}

		var claims types.AccessTokenClaims
		var err error
		if claims, err = utils.ValidateJWT(token, jwtCfg.GetSecret(), lgr); err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [RefreshTokenHandler] [ValidateJWT] %s", err.Error()))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
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
			lgr.Debug(fmt.Sprintf("[Handlers] [RefreshTokenHandler] [IssueJWT] %s", err.Error()))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: tkn,
			RefreshToken:        token,
		}
		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func ValidateTokenHandler(lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// If claims are present on context then tokens are already validated by auth middleware
		token := req.Context().Value("claims")

		if token == nil {
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusUnauthorized, "claims are invalid"), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, "claims are valid", w, lgr)
	}
}
