package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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
			switch errx.Kind() {
			case custom_errors.DuplicateRecordInsertion:
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, errx.Error()), w, lgr)
			default:
				utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.Error()), w, lgr)
			}
			return
		}

		resp := types.CreateFolderResponse{
			Name:     data.Name,
			FolderID: newFolderID,
		}
		utils.WriteSuccessResponse(http.StatusOK, resp, w, lgr)
	}
}

func GetFolderHandler(svc service.FoldersService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)
		paramsMap := req.Context().Value("url_params").(map[string]string)

		if paramsMap["folderID"] == "" {
			lgr.Info("[Handlers] [GetFolderHandler] folderID URL parameter empty")
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "folderID parameter not specified"), w, lgr)
			return
		}

		id, err := strconv.Atoi(paramsMap["folderID"])
		if err != nil {
			lgr.Info("[Handlers] [GetFolderHandler] [Atoi] specified folderID is not a number")
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "specified folderID is of incorrect type"), w, lgr)
			return
		}

		folderContents, errx := svc.Get(types.FolderID(id), claims)
		if errx != nil {
			errMsg := fmt.Sprintf("[Handlers] [GetFolderHandler] [Get] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, folderContents, w, lgr)
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
