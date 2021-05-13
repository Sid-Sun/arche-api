package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/contract"
	"github.com/sid-sun/arche-api/app/http/resperr"
	"go.uber.org/zap"
)

func writeResponse(code int, data []byte, resp http.ResponseWriter, lgr *zap.Logger) {
	resp.WriteHeader(code)
	_, err := resp.Write(data)
	if err != nil {
		lgr.Debug(fmt.Sprintf("[Utils] [ResponseWriter] [writeResponse] [Write] %s", err.Error()))
		erx.WithArgs(err, erx.SeverityDebug, erx.Kind("ResponseWriteError"))
	}
}

func writeAPIResponse(code int, ar contract.APIResponse, resp http.ResponseWriter, lgr *zap.Logger) {
	b, err := json.Marshal(&ar)
	if err != nil {
		lgr.Error(fmt.Sprintf("[Utils] [ResponseWriter] [writeAPIResponse] [Marshal] %s", err.Error()))
		writeResponse(http.StatusInternalServerError, []byte("internal server error"), resp, lgr)
		return
	}

	writeResponse(code, b, resp, lgr)
}

func WriteSuccessResponse(statusCode int, data interface{}, resp http.ResponseWriter, lgr *zap.Logger) {
	writeAPIResponse(statusCode, contract.NewSuccessResponse(data), resp, lgr)
}

func WriteFailureResponse(gr resperr.ResponseError, resp http.ResponseWriter, lgr *zap.Logger) {
	writeAPIResponse(gr.StatusCode(), contract.NewFailureResponse(gr.Description()), resp, lgr)
}
