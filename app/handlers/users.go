package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/http/resperr"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

func CreateUserHandler(svc service.UsersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "could not read request body"), w, lgr)
			return
		}

		var data types.UserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}

		key, hash, err := utils.GenerateEncryptionKey(lgr)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [GenerateEncryptionKey] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, "could not generate encryption keys"), w, lgr)
			return
		}

		var decryptedKey []byte
		copy(decryptedKey, key)
		utils.EncryptKey(key, data.Password, lgr)

		usr, errx := svc.CreateUser(data.Email, key, hash)
		if errx != nil {
			if errx.Kind() == custom_errors.DuplicateRecordInsertion {
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "user already exists"), w, lgr)
				return
			}
			switch errx.Severity {
			case erx.SeverityError:
				lgr.Error(fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [CreateUser] %s", errx.Error()))
			case erx.SeverityInfo:
				lgr.Info(fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [CreateUser] %s", errx.Error()))
			default:
				lgr.Debug(fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [CreateUser] %s", errx.Error()))
			}
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		accessToken, refreshToken, err := utils.IssueTokens(usr.ID, decryptedKey, cfg, lgr)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [IssueTokens] %s", errx.Error()))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: accessToken,
			RefreshToken:        refreshToken,
		}

		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func LoginUserHandler(svc service.UsersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "could not read request body"), w, lgr)
			return
		}

		var data types.UserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}

		usr, errx := svc.GetUser(data.Email)
		if errx != nil {
			if errx.Kind() == custom_errors.NoRowsInResultSet {
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "user does not exist"), w, lgr)
				return
			}
			switch errx.Severity {
			case erx.SeverityError:
				lgr.Error(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [GetUser] %s", errx.Error()))
			case erx.SeverityInfo:
				lgr.Info(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [GetUser] %s", errx.Error()))
			default:
				lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [GetUser] %s", errx.Error()))
			}
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		key, err := base64.StdEncoding.DecodeString(usr.EncryptionKey)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [DecodeString] EncryptionKey %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		_ = utils.DecryptKey(key, data.Password, lgr)
		hash, err := base64.StdEncoding.DecodeString(usr.KeyHash)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [DecodeString] Hash %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		keyHash := sha3.Sum256(key)
		if !bytes.Equal(hash, keyHash[:]) {
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusUnauthorized, "incorrect password"), w, lgr)
			return
		}

		accessToken, refreshToken, err := utils.IssueTokens(usr.ID, key, cfg, lgr)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [DecodeString] Hash %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: accessToken,
			RefreshToken:        refreshToken,
		}

		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}
