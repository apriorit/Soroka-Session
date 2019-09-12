package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	er "github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	multierror "github.com/hashicorp/go-multierror"
)

func TestMultiError(t *testing.T) {
	err := er.ErrNonAuthorized
	err = multierror.Append(err, er.ErrExpiredAccessToken)

	_, ok := err.(*multierror.Error)
	assert.True(t, ok)
}

func TestFormUnixData(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	expire1 := time.Unix(time.Date(1976, 1, 1, 0, 0, 0, 0, loc).Unix(), 0).In(loc)
	expire2 := time.Unix(time.Now().Add(1).Unix(), 0).In(loc)

	fmt.Println(expire1)
	fmt.Println(expire2)
	assert.Equal(t, loc.String(), "UTC")
}

func TestGetCookieWithExpiredToken(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	expire := time.Unix(time.Date(1976, 1, 1, 0, 0, 0, 0, loc).Unix(), 0).In(loc)

	testCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}

	cookie := GetCookieWithExpiredToken(nil, loc)

	assert.Equal(t, testCookie, cookie)
}

func TestGetCookieWithNewToken(t *testing.T) {
	expireRaw := time.Now().Add(time.Duration(720) * time.Hour).Unix()
	loc, _ := time.LoadLocation("UTC")
	expireUnixTime := time.Unix(expireRaw, 0).In(loc)

	tokenData := m.TokenData{
		Token:          "x78Fxkjk=",
		Type:           "Bearer",
		ExpirationDate: expireRaw,
	}

	testCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "x78Fxkjk=",
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expireUnixTime,
		Secure:   true,
		HttpOnly: true,
	}

	cookie := GetCookieWithNewToken(&tokenData, loc)

	assert.Equal(t, testCookie, cookie)
}
