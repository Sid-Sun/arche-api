package handlers

import (
	"encoding/json"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func CreateFolderHandler(svc service.FoldersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("authentication_token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var claims types.AccessTokenClaims
		var err error
		if claims, err = utils.ValidateJWT(token, cfg.GetSecret(), lgr); err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data types.CreateFolderRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newFolderID, err := svc.Create(data.Name, claims)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := types.CreateFolderResponse{
			Name:     data.Name,
			FolderID: newFolderID,
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

func GetFoldersHandler(svc service.FoldersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("authentication_token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var claims types.AccessTokenClaims
		var err error
		if claims, err = utils.ValidateJWT(token, cfg.GetSecret(), lgr); err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		folders, err := svc.GetAll(claims)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(folders)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			// TODO: Add Logging
			return
		}
	}
}

func DeleteFolder(svc service.FoldersService, cfg *config.JWTConfig, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("authentication_token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var claims types.AccessTokenClaims
		var err error
		if claims, err = utils.ValidateJWT(token, cfg.GetSecret(), lgr); err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var data types.DeleteFolderRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = svc.Delete(data.FolderID, claims)
		if err != nil {
			// TODO: Add Logging
			// TODO: Check if a folder was deleted, don't allow non-existent folders to be deleted
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, err = w.Write(d)
		if err != nil {
			// TODO: Add Logging
			return
		}
	}
}
