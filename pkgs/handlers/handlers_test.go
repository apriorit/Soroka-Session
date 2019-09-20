package handlers

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	"github.com/stretchr/testify/assert"
)

func TestDecodeLoginRequest(t *testing.T) {
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/sessions/login", nil)
	assert.NoError(t, err)

	rawRequest.SetBasicAuth("admin", "a@m1n")

	testData := endpoints.LoginRequest{
		Req: models.LoginData{
			UserName: "admin",
			Password: "a@m1n",
		},
	}

	resp, err := DecodeLoginRequest(context.Background(), rawRequest)

	assert.Equal(t, testData, resp)
}

func TestDecodeLogoutRequest(t *testing.T) {
	rawRequest, err := http.NewRequest("GET", "https://edms.com/api/v1/sessions/logout", nil)
	assert.NoError(t, err)

	cookie := &http.Cookie{
		Name:  "refresh_token",
		Value: "x78H56Bar90=",
	}

	rawRequest.AddCookie(cookie)

	testData := endpoints.LogoutRequest{
		Req: models.LogoutData{
			Cookie: cookie,
		},
	}

	resp, err := DecodeLogoutRequest(context.Background(), rawRequest)
	assert.NoError(t, err)
	req, ok := resp.(endpoints.LogoutRequest)
	assert.True(t, ok)
	assert.Equal(t, testData.Req.Cookie.Name, req.Req.Cookie.Name)
	assert.Equal(t, testData.Req.Cookie.Value, req.Req.Cookie.Value)
}

func TestDecodeCheckTokenRequest(t *testing.T) {
	var rawTokenString = []byte(`{"access_token": "x78H56Bar90="}`)
	rawRequest, err := http.NewRequest("POST", "https://edms.com/api/v1/users/check_token", bytes.NewBuffer(rawTokenString))
	assert.NoError(t, err)

	cookie := &http.Cookie{
		Name:  "refresh_token",
		Value: "x34H56Bar45=",
	}

	rawRequest.AddCookie(cookie)

	testData := endpoints.CheckTokenRequest{
		Req: models.CheckTokenServiceInput{
			AccessToken:  "x78H56Bar90=",
			RefreshToken: "x34H56Bar45=",
		},
	}

	resp, err := DecodeCheckTokenRequest(context.Background(), rawRequest)
	assert.NoError(t, err)
	req, ok := resp.(endpoints.CheckTokenRequest)
	assert.True(t, ok)
	assert.Equal(t, req.Req.AccessToken, testData.Req.AccessToken)
	assert.Equal(t, req.Req.RefreshToken, testData.Req.RefreshToken)
}
