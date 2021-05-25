package handlers

import (
	"errors"
	"regexp"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
)

func validateEmail(email string) *erx.Erx {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if ok := emailRegex.MatchString(email); !ok {
		return erx.WithArgs(errors.New("invalid email syntax"), erx.SeverityInfo, custom_errors.InvalidEmailAddress)
	}
	return nil
}
