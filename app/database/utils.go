package database

import (
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/nsnikhil/erx"
)

func checkForSQLError(err error) (*mssql.Error, *erx.Erx) {
	var errx *erx.Erx
	errx = erx.WithArgs(err, erx.SeverityDebug)
	if mssqlError, ok := err.(mssql.Error); ok {
		errx = erx.WithArgs(mssqlError, erx.SeverityDebug)
		return &mssqlError, errx
	}
	return nil, errx
}
