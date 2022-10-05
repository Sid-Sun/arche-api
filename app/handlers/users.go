package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/http/resperr"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

func CreateUserHandler(svc service.UsersService, veCfg *config.VerificationEmailConfig, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "could not read request body"), w, lgr)
			return
		}

		var data types.CreateUserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}
		data.Email = strings.ToLower(data.Email)

		if errx := validateEmail(data.Email); errx != nil {
			lgr.Info(fmt.Sprintf("[Handlers] [Users] [validateEmail] [InvalidEmail] %v", errx.String()))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "email address is not valid"), w, lgr)
			return
		}

		key, hash, err := utils.GenerateEncryptionKey(lgr)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateUserHandler] [GenerateEncryptionKey] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, "could not generate encryption keys"), w, lgr)
			return
		}

		decryptedKey := make([]byte, len(key))
		copy(decryptedKey, key)
		utils.EncryptKey(key, data.Password, lgr)
		verificationToken := utils.RandString(veCfg.GetTokenLength())

		_, errx := svc.CreateUser(data.Email, key, hash, verificationToken)
		if errx != nil {
			if errx.Kind() == custom_errors.DuplicateRecordInsertion {
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "user already exists"), w, lgr)
				return
			}
			errMsg := fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [CreateUser] %s", errx.Error())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		resp := types.CreateUserResponse{
			VerificationEmailSent: false,
			UserCreated:           true,
		}
		fmt.Printf("here")

		errx = svc.SendVerificationEmail(data.Email, verificationToken, data.VerificationCallbackURL, veCfg)
		if errx != nil {
			fmt.Printf(errx.String())
			errMsg := fmt.Sprintf("[Handlers] [Users] [CreateUserHandler] [SendVerificationEmail] %s", errx.Error())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
			return
		}
		resp.VerificationEmailSent = true
		fmt.Printf("where")

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

		var data types.LoginUserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}
		data.Email = strings.ToLower(data.Email)

		if errx := validateEmail(data.Email); errx != nil {
			lgr.Info(fmt.Sprintf("[Handlers] [Users] [validateEmail] [InvalidEmail] %v", errx.String()))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "email address is not valid"), w, lgr)
			return
		}

		usr, errx := svc.GetUser(data.Email)
		if errx != nil {
			if errx.Kind() == custom_errors.NoRowsInResultSet {
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "user does not exist"), w, lgr)
				return
			}
			errMsg := fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [GetUser] %s", errx.Error())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
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

		resp := types.LoginUserResponse{
			VerificationPending: true,
		}

		if !usr.Verified {
			lgr.Info("[Handlers] [Users] [LoginUserHandler] [VerifiedCheck] User is not verified")
			utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
			return
		}

		resp.VerificationPending = false
		resp.AuthenticationToken, resp.RefreshToken, err = utils.IssueTokens(usr.ID, key, cfg, lgr)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [LoginUserHandler] [IssueTokens] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func ActivateUserHandler(svc service.UsersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "could not read request body"), w, lgr)
			return
		}

		var data types.ActivateUserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}

		errx := svc.ActivateUser(data.VerificationToken)
		if errx != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [ActivateUser] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, err.Error()), w, lgr)
			return
		}

		resp := types.ActivateUserResponse{
			VerificationPending: false,
		}

		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func ResendValidationHandler(svc service.UsersService, veCfg *config.VerificationEmailConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "could not read request body"), w, lgr)
			return
		}

		var data types.ResendVerificationRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "an error occoured when unmarshaling JSON"), w, lgr)
			return
		}
		data.Email = strings.ToLower(data.Email)

		verified, token, errx := svc.GetVerificationStatus(data.Email)
		if errx != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [ActivateUser] %v", errx))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, errx.Error()), w, lgr)
			return
		}

		if verified {
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "user is already verified"), w, lgr)
			return
		}

		if token == "" {
			token = utils.RandString(veCfg.GetTokenLength())
			errx = svc.UpdateVerificationToken(data.Email, token)
			if errx != nil {
				lgr.Debug(fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [UpdateVerificationToken] %v", errx))
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.Error()), w, lgr)
				return
			}
		}

		errx = svc.SendVerificationEmail(data.Email, token, data.VerificationCallbackURL, veCfg)
		if errx != nil {
			fmt.Println(errx.String())
			errMsg := fmt.Sprintf("[Handlers] [Users] [ActivateUserHandler] [SendVerificationEmail] %s", errx.Error())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.Error()), w, lgr)
			return
		}

		resp := types.ResendVerificationResponse{
			VerificationEmailSent: true,
		}

		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}
