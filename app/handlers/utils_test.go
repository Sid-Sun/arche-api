package handlers

import (
	"testing"

	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	testEmail := "jane@example.com"
	errx := validateEmail(testEmail)
	assert.Nil(t, errx)

	testEmail = "john@game.com"
	errx = validateEmail(testEmail)
	assert.Nil(t, errx)

	testEmail = "_@lol.com"
	errx = validateEmail(testEmail)
	assert.Nil(t, errx)

	testEmail = `5=[43]\@apple.com`
	errx = validateEmail(testEmail)
	assert.NotNil(t, errx)
	assert.Equal(t, custom_errors.InvalidEmailAddress, errx.Kind())

	testEmail = `-=\@t-=oast.=com`
	errx = validateEmail(testEmail)
	assert.NotNil(t, errx)
	assert.Equal(t, custom_errors.InvalidEmailAddress, errx.Kind())

	testEmail = "potato"
	errx = validateEmail(testEmail)
	assert.NotNil(t, errx)
	assert.Equal(t, custom_errors.InvalidEmailAddress, errx.Kind())

	testEmail = "potato@eample"
	errx = validateEmail(testEmail)
	assert.NotNil(t, errx)
	assert.Equal(t, custom_errors.InvalidEmailAddress, errx.Kind())
}
