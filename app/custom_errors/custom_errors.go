package custom_errors

import (
	"errors"
	"github.com/nsnikhil/erx"
)

const DuplicateRecordInsertion = erx.Kind("DuplicateRecordInsertion")
const NoRowsInResultSet = erx.Kind("NoRowsInResultSet")
const NoRowsAffected = erx.Kind("NoRowsAffected")

var ErrSQLNoResultsInSet = errors.New("sql: no rows in result set")
