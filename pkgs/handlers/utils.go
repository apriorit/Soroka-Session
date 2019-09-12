package handlers

import (
	"net/http"
	"time"

	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//AddCookie adds cookie http-only in response
func AddCookie(w http.ResponseWriter, op string, rt *m.TokenData) {
	var cookie http.Cookie
	loc, _ := time.LoadLocation("UTC")

	if op == "new" {
		cookie = GetCookieWithNewToken(rt, loc)
	} else {
		cookie = GetCookieWithExpiredToken(nil, loc)
	}

	http.SetCookie(w, &cookie)
}

//GetCookieWithNewToken forms cookie with a new refresh token
func GetCookieWithNewToken(rt *m.TokenData, loc *time.Location) http.Cookie {
	expire := time.Unix(rt.ExpirationDate, 0).In(loc)

	return http.Cookie{
		Name:     "refresh_token",
		Value:    rt.Token,
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}
}

//GetCookieWithExpiredToken forms cookie with expired refresh token
func GetCookieWithExpiredToken(rt *m.TokenData, loc *time.Location) http.Cookie {
	expire := time.Unix(time.Date(1976, 1, 1, 0, 0, 0, 0, loc).Unix(), 0).In(loc)

	return http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "edms.com",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}
}
