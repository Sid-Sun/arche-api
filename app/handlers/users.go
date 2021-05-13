package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"net/http"
	"time"
)

func CreateUserHandler(svc service.UsersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [CreateUserHandler] [ReadAll] %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data types.UserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [CreateUserHandler] [Unmarshal] %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key, hash, err := utils.GenerateEncryptionKey(lgr)
		if err != nil {
			// TODO: Add Logging
			return
		}
		err = utils.EncryptKey(key, data.Password, lgr)
		usr, err := svc.CreateUser(data.Email, key, hash)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims := types.AccessTokenClaims{
			UserID:        usr.ID,
			EncryptionKey: key,
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(cfg.GetTTL())).Unix(),
			},
		}
		tkn, err := utils.IssueJWT(claims, cfg.GetSecret(), lgr)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: tkn,
		}

		d, err = json.Marshal(resp)
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

func LoginUserHandler(svc service.UsersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [LoginUserHandler] [ReadAll] %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data types.UserRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [LoginUserHandler] [Unmarshal] %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		usr, err := svc.GetUser(data.Email)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [LoginUserHandler] [GetUser] %v", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		key, err := base64.StdEncoding.DecodeString(usr.EncryptionKey)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [LoginUserHandler] [DecodeString] EncryptionKey %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = utils.DecryptKey(key, data.Password, lgr)
		hash, err := base64.StdEncoding.DecodeString(usr.KeyHash)
		if err != nil {
			lgr.Error(fmt.Sprintf("[Handlers] [LoginUserHandler] [DecodeString] Hash %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		keyHash := sha3.Sum256(key)
		if !bytes.Equal(hash, keyHash[:]) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := types.AccessTokenClaims{
			UserID:        usr.ID,
			EncryptionKey: key,
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}
		tkn, err := utils.IssueJWT(claims, cfg.GetSecret(), lgr)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := types.UserResponse{
			AuthenticationToken: tkn,
		}

		d, err = json.Marshal(resp)
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
