package custom_errors

import (
	"github.com/nsnikhil/erx"
)

const DuplicateRecordInsertion = erx.Kind("DuplicateRecordInsertion")
const NoRowsInResultSet = erx.Kind("NoRowsInResultSet")
const NoRowsAffected = erx.Kind("NoRowsAffected")
const InvalidEmailAddress = erx.Kind("InvalidEmailAddress")
