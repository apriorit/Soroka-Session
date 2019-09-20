package handlers

import (
	"net/http"
	"time"
)

//AddCookie adds cookie http-only in response
func AddCookie(w http.ResponseWriter, tokenValue string, expiresIn int64) {
	cookie := GetCookieWithToken(tokenValue, expiresIn)
	http.SetCookie(w, &cookie)
}

func GetCookieWithToken(value string, expiresIn int64) http.Cookie {
	loc, _ := time.LoadLocation("UTC")
	expire := time.Unix(expiresIn, 0).In(loc)

	return http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}
}
