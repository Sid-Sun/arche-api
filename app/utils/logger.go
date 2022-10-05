package utils

import (
	"github.com/nsnikhil/erx"
	"go.uber.org/zap"
)

func LogWithSeverity(text string, severity erx.Severity, lgr *zap.Logger) {
	switch severity {
	case erx.SeverityError:
		lgr.Error(text)
	case erx.SeverityInfo:
		lgr.Info(text)
	case erx.SeverityWarn:
		lgr.Warn(text)
	case erx.SeverityDebug:
	default:
		lgr.Debug(text)
	}
}
