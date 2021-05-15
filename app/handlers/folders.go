package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/http/resperr"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

// TODO: Error Handling from DB (errx)
// TODO: Severity-Checked Logging for errx

func CreateFolderHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateFolderHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "failed to read request body"), w, lgr)
			return
		}

		var data types.CreateFolderRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateFolderHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
			return
		}

		newFolderID, errx := svc.Create(data.Name, claims)
		if errx != nil {
			errMsg := fmt.Sprintf("[Handlers] [CreateFolderHandler] [Create] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		resp := types.CreateFolderResponse{
			Name:     data.Name,
			FolderID: newFolderID,
		}
		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func GetFoldersHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		folders, errx := svc.GetAll(claims)
		if errx != nil {
			errMsg := fmt.Sprintf("[Handlers] [GetFoldersHandler] [GetAll] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, folders, w, lgr)
	}
}

func DeleteFolderHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [DeleteFolderHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "failed to read request body"), w, lgr)
			return
		}

		var data types.DeleteFolderRequest
		err = json.Unmarshal(d, &data)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [DeleteFolderHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
			return
		}

		errx := svc.Delete(data.FolderID, claims)
		if errx != nil && errx.Kind() != custom_errors.NoRowsAffected {
			errMsg := fmt.Sprintf("[Handlers] [DeleteFolderHandler] [GetAll] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		if errx.Kind() == custom_errors.NoRowsAffected {
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "folder does not exist or doesn't belong to user"), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, data.FolderID, w, lgr)
	}
}
