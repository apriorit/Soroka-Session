package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	"github.com/stretchr/testify/assert"
)

func TestGetCookieWithExpiredToken(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	expire := time.Unix(0, 0).In(loc)

	testCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}

	cookie := GetCookieWithToken("", 0)

	assert.Equal(t, testCookie, cookie)
}

func TestGetCookieWithNewToken(t *testing.T) {
	expireRaw := time.Now().Add(time.Duration(720) * time.Hour).Unix()
	loc, _ := time.LoadLocation("UTC")
	expireUnixTime := time.Unix(expireRaw, 0).In(loc)

	tokenData := models.TokenData{
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

	cookie := GetCookieWithToken(tokenData.Token, tokenData.ExpirationDate)

	assert.Equal(t, testCookie, cookie)
}
