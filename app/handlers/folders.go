package handlers

import (
	"encoding/json"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func CreateFolderHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

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

func GetFoldersHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

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

func DeleteFolder(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

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
