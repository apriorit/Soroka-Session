package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/db"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type InputType int

const (
	AccessExpiredRefreshNotExpired InputType = iota
	AccessNotExpiredRefreshNotExpired
	RefreshExpired
	InvalidToken
)

func PrepareCookie(t string) (*http.Cookie, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return &http.Cookie{}, err
	}
	tokenData := models.TokenData{
		Token:          t,
		Type:           "Bearer",
		ExpirationDate: time.Now().Add(time.Duration(720) * time.Hour).Unix(),
	}

	return &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenData.Token,
		Path:     "/",
		Domain:   "edms.com",
		Expires:  time.Unix(tokenData.ExpirationDate, 0).In(loc),
		Secure:   true,
		HttpOnly: true,
	}, nil
}

func PrepareServiceAndDb(id, token string) models.ISessionService {
	db, _ := db.Connection(config.GetLogger().Logger, "stub")
	var testData = struct {
		sub   string
		token string
	}{
		id,
		token,
	}
	db.Save(testData.sub, testData.token)
	return Build(log.NewNopLogger(), db, []byte("secret"), make([]byte, 0))
}

func PrepareLogoutRequest(c *http.Cookie) *models.LogoutData {
	return &models.LogoutData{
		Cookie: c,
	}
}

func CreateAccessToken(iss, cid string, mask int64, expired bool) (tokenString string, err error) {
	var exp int64
	iat := time.Now().Unix()

	if expired {
		exp = time.Now().Add(time.Duration(0) * time.Second).Unix()
	} else {
		exp = time.Now().Add(time.Duration(876000) * time.Hour).Unix()
	}

	claims := jwt.MapClaims{
		"iss":  iss,
		"sub":  cid,
		"iat":  iat,
		"nbf":  iat,
		"exp":  exp,
		"mask": mask,
		"aud":  []string{cid, iss},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err = token.SignedString([]byte("secret")); err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(iss, cid string, expired bool) (tokenString string, err error) {
	var exp int64
	iat := time.Now().Unix()

	if expired {
		exp = time.Now().Add(time.Duration(0) * time.Second).Unix()
	} else {
		exp = time.Now().Add(time.Duration(876000) * time.Hour).Unix()
	}

	claims := jwt.MapClaims{
		"iss": iss,
		"sub": cid,
		"iat": iat,
		"nbf": iat,
		"exp": exp,
		"aud": []string{cid, iss},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err = token.SignedString([]byte("secret")); err != nil {
		return "", err
	}

	return tokenString, nil
}

func PrepareCheckTokenInput(t InputType) (input models.CheckTokenServiceInput, err error) {
	switch t {
	case AccessExpiredRefreshNotExpired:
		if input.AccessToken, err = CreateAccessToken(constants.TokenIssuer, "gladys.champl@edms.com", 32767, true); err != nil {
			return input, err
		}
		if input.RefreshToken, err = CreateRefreshToken(constants.TokenIssuer, "gladys.champl@edms.com", false); err != nil {
			return input, err
		}
	case AccessNotExpiredRefreshNotExpired:
		if input.AccessToken, err = CreateAccessToken(constants.TokenIssuer, "gladys.champl@edms.com", 32767, false); err != nil {
			return input, err
		}
		if input.RefreshToken, err = CreateRefreshToken(constants.TokenIssuer, "gladys.champl@edms.com", false); err != nil {
			return input, err
		}
	case RefreshExpired:
		if input.AccessToken, err = CreateAccessToken(constants.TokenIssuer, "gladys.champl@edms.com", 32767, true); err != nil {
			return input, err
		}
		if input.RefreshToken, err = CreateRefreshToken(constants.TokenIssuer, "gladys.champl@edms.com", true); err != nil {
			return input, err
		}
	case InvalidToken:
		input.AccessToken = "Z2xhZHlzLmNoYW1wbEBlZG1zLmNvbTphQG0xbg"
		input.RefreshToken = "Td5MZHlzHc4PYW1wbsdsZG1zLmNvbTphQG0d6f"
	}

	return input, nil
}

func TestLogout_ValidUser(t *testing.T) {
	userID := "gladys.champl@edms.com"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZ2xhZHlzLmNoYW1wbEBlZG1zLmNvbSIsImh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiXSwiZXhwIjoxNTY4Mzc4NTg0LCJpYXQiOjE1NjgyOTIxODQsImlzcyI6Imh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiLCJtYXNrIjozMjc2NywibmJmIjoxNTY4MjkyMTg0LCJzdWIiOiJnbGFkeXMuY2hhbXBsQGVkbXMuY29tIn0.xUKkiNOClwhhGlgXWj_J9u0t_ImJKsW-mbK9xuTiF5o"
	cookie, err := PrepareCookie(token)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb(userID, token)
	req := PrepareLogoutRequest(cookie)

	svc.Logout(context.Background(), *req)
	assert.NoError(t, err)
}

func TestLogout_InvalidUser(t *testing.T) {
	userID := "user@email.com"
	token := "Z2xhZHlzLmNoYW1wbEBlZG1zLmNvbTphQG0xbg"
	cookie, err := PrepareCookie(token)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb(userID, token)
	req := PrepareLogoutRequest(cookie)

	err = svc.Logout(context.Background(), *req)
	assert.Error(t, err)
}

func TestCheckToken_ExpiredAccessToken(t *testing.T) {
	input, err := PrepareCheckTokenInput(AccessExpiredRefreshNotExpired)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb("gladys.champl@edms.com", input.RefreshToken)
	//Wait until the token expiration date
	time.Sleep(1 * time.Second)
	output, err := svc.CheckToken(context.Background(), input)
	assert.NoError(t, err)
	//New access token expected
	assert.True(t, input.AccessToken != output.AccessToken)
}

func TestCheckToken_NotExpiredAccessToken(t *testing.T) {
	input, err := PrepareCheckTokenInput(AccessNotExpiredRefreshNotExpired)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb("gladys.champl@edms.com", input.RefreshToken)
	output, err := svc.CheckToken(context.Background(), input)
	assert.NoError(t, err)
	//New access token expected
	assert.True(t, input.AccessToken == output.AccessToken)
}

func TestCheckToken_ExpiredRefreshToken(t *testing.T) {
	input, err := PrepareCheckTokenInput(RefreshExpired)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb("gladys.champl@edms.com", input.RefreshToken)
	//Wait until the token expiration date
	time.Sleep(1 * time.Second)
	output, err := svc.CheckToken(context.Background(), input)
	assert.Error(t, err)
	//New access token expected
	assert.True(t, "" == output.AccessToken)
}

func TestCheckToken_InvalidToken(t *testing.T) {
	input, err := PrepareCheckTokenInput(InvalidToken)
	assert.NoError(t, err)
	svc := PrepareServiceAndDb("gladys.champl@edms.com", input.RefreshToken)
	output, err := svc.CheckToken(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, output.AccessToken, "")
}

func TestCheckToken_NoSuchTOkenInDb(t *testing.T) {
	input, err := PrepareCheckTokenInput(AccessExpiredRefreshNotExpired)
	assert.NoError(t, err)
	//Wait until the token expiration date
	time.Sleep(2 * time.Second)
	svc := PrepareServiceAndDb("user@email.com", input.RefreshToken)
	output, err := svc.CheckToken(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, output.AccessToken, "")
}
