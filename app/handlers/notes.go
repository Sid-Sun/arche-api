package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sid-sun/arche-api/app/http/resperr"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

func GetNotesHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		notes, errx := svc.GetAll(claims)
		if errx != nil {
			errMsg := fmt.Sprintf("[Handlers] [GetNotesHandler] [GetAll] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, notes, w, lgr)
	}
}

func GetNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)
		paramsMap := req.Context().Value("url_params").(map[string]string)

		if paramsMap["noteID"] == "" {
			lgr.Info("[Handlers] [GetNoteHandler] noteID URL parameter empty")
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "noteID parameter not specified"), w, lgr)
			return
		}

		id, err := strconv.Atoi(paramsMap["noteID"])
		if err != nil {
			lgr.Info("[Handlers] [GetNoteHandler] [Atoi] specified noteID is not a number")
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "specified noteID is of incorrect type"), w, lgr)
			return
		}

		notes, errx := svc.Get(types.NoteID(id), claims)
		if errx != nil {
			errMsg := fmt.Sprintf("[Handlers] [GetNoteHandler] [Get] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, notes, w, lgr)
	}
}

func CreateNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateNoteHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "failed to read request body"), w, lgr)
			return
		}

		var body types.CreateNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [CreateNoteHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
			return
		}

		noteID, errx := svc.Create(body.Name, body.Data, body.FolderID, claims)
		if errx != nil {
			// TODO: Implement Non-Existent Data operation or Unauthorized data operation errors
			errMsg := fmt.Sprintf("[Handlers] [GetNotesHandler] [Create] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		res := types.CreateNoteResponse{
			NoteID: noteID,
		}
		utils.WriteSuccessResponse(http.StatusOK, res, w, lgr)
	}
}

func UpdateNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [UpdateNoteHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "failed to read request body"), w, lgr)
			return
		}

		var body types.UpdateNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [UpdateNoteHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
			return
		}

		errx := svc.Update(body.Name, body.Data, body.FolderID, body.NoteID, claims)
		if errx != nil {
			// TODO: Implement Non-Existent Data operation or Unauthorized data operation errors
			errMsg := fmt.Sprintf("[Handlers] [UpdateNoteHandler] [Update] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, types.UpdateNoteResponse(body), w, lgr)
	}
}

func DeleteNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [DeleteNoteHandler] [ReadAll] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, "failed to read request body"), w, lgr)
			return
		}

		var body types.DeleteNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			lgr.Debug(fmt.Sprintf("[Handlers] [DeleteNoteHandler] [Unmarshal] %v", err))
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusBadRequest, err.Error()), w, lgr)
			return
		}

		errx := svc.Delete(body.NoteID, claims)
		if errx != nil {
			// TODO: Implement Non-Existent Data operation or Unauthorized data operation errors
			errMsg := fmt.Sprintf("[Handlers] [DeleteNoteHandler] [Delete] %v", errx.String())
			utils.LogWithSeverity(errMsg, errx.Severity, lgr)
			utils.WriteFailureResponse(resperr.NewResponseError(http.StatusInternalServerError, errx.String()), w, lgr)
			return
		}

		utils.WriteSuccessResponse(http.StatusOK, types.DeleteNoteResponse(body), w, lgr)
	}
}
