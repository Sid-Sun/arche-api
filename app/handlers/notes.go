package handlers

import (
	"encoding/json"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func GetNotesHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

		notes, err := svc.GetAll(claims)
		if err != nil {
			// TODO: Add Logging
			return
		}

		data, err := json.Marshal(notes)
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

func CreateNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var body types.CreateNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		noteID, err := svc.Create(body.Name, body.Data, body.FolderID, claims)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res := types.CreateNoteResponse{
			NoteID: noteID,
		}

		d, err = json.Marshal(res)
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

func UpdateNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var body types.UpdateNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = svc.Update(body.Name, body.Data, body.FolderID, body.NoteID, claims)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func DeleteNoteHandler(svc service.NotesService, lgr *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var claims types.AccessTokenClaims
		claims = req.Context().Value("claims").(types.AccessTokenClaims)

		d, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var body types.DeleteNoteRequest
		err = json.Unmarshal(d, &body)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = svc.Delete(body.NoteID, claims)
		if err != nil {
			// TODO: Add Logging
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
