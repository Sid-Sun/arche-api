package service

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/types"
	"github.com/sid-sun/arche-api/app/utils"
	"go.uber.org/zap"
)

func decryptNote(note types.Note, blockCipher cipher.Block, lgr *zap.Logger) (types.Note, *erx.Erx) {
	name, err := base64.StdEncoding.DecodeString(note.Name)
	if err != nil {
		lgr.Debug(fmt.Sprintf("[Service] [Utils] [decryptNote] [DecodeString] [Name] %s", err.Error()))
		return types.Note{}, erx.WithArgs(err, erx.SeverityDebug)
	}
	note.Name = string(utils.CFBDecrypt(name, blockCipher))

	data, err := base64.StdEncoding.DecodeString(note.Data)
	if err != nil {
		lgr.Debug(fmt.Sprintf("[Service] [Utils] [decryptNote] [DecodeString] [Data] %s", err.Error()))
		return types.Note{}, erx.WithArgs(err, erx.SeverityDebug)
	}
	note.Data = string(utils.CFBDecrypt(data, blockCipher))
	return note, nil
}
