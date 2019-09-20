package service

import (
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
)

func TestCreatePayload_Access(t *testing.T) {
	testPayload, exp, err := CreatePayload(access, "user@example.com", constants.TokenIssuer, 2048)
	assert.NoError(t, err)
	assert.NoError(t, testPayload.Valid())
	assert.True(t, testPayload.VerifyIssuer(constants.TokenIssuer, true))
	assert.NotZero(t, exp)
}

func TestCreatePayload_Refresh(t *testing.T) {
	testPayload, exp, err := CreatePayload(refresh, "user@example.com", constants.TokenIssuer, 2048)
	assert.NoError(t, err)
	assert.NoError(t, testPayload.Valid())
	assert.True(t, testPayload.VerifyIssuer(constants.TokenIssuer, true))
	assert.NotZero(t, exp)
}

func TestGenerateToken_Access(t *testing.T) {
	sStub := SessionsService{
		client: &http.Client{},
		secret: []byte("secret"),
		Logger: config.GetLogger().Logger,
	}

	var testData = struct {
		sub  string
		iss  string
		aud  []string
		mask int64
	}{
		"user@example.com",
		"https://edms.com/sessions",
		[]string{"user@example.com", "https://edms.com/sessions"},
		2048,
	}

	tokenServiceData, err := sStub.GenerateToken(access, "user@example.com", 2048)
	assert.NoError(t, err)
	assert.NotEqual(t, tokenServiceData.Token, "")

	//Check token validness and token claims
	claims := jwt.MapClaims{}
	rawToken, err := jwt.ParseWithClaims(tokenServiceData.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	assert.True(t, rawToken.Valid)

	//Check nbf, exp and iat claims
	assert.NoError(t, claims.Valid())
	//Check the number of claims
	assert.Equal(t, len(claims), constants.ClaimsPerAccessToken)

	maskClaim, ok := claims["mask"].(float64)
	assert.True(t, ok)

	//Check other claims
	assert.Equal(t, testData.sub, claims["sub"])
	assert.Equal(t, testData.iss, claims["iss"])
	assert.Equal(t, testData.mask, int64(maskClaim))
}

func TestGenerateToken_Refresh(t *testing.T) {
	sStub := SessionsService{
		client: &http.Client{},
		secret: []byte("secret"),
		Logger: config.GetLogger().Logger,
	}

	var testData = struct {
		sub string
		iss string
		aud []string
	}{
		"user@example.com",
		"https://edms.com/sessions",
		[]string{"user@example.com", "https://edms.com/sessions"},
	}

	token, err := sStub.GenerateToken(refresh, "user@example.com", 2048)
	assert.NoError(t, err)
	assert.NotEqual(t, token.Token, "")

	//Check token claims
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	//Check nbf, exp and iat claims
	assert.NoError(t, claims.Valid())
	//Check the number of claims
	assert.Equal(t, len(claims), constants.ClaimsPerRefreshToken)

	//Check other claims
	assert.Equal(t, testData.sub, claims["sub"])
	assert.Equal(t, testData.iss, claims["iss"])
}

func TestCheckTokenValidness_ValidTokenValidClaim(t *testing.T) {
	sStub := SessionsService{
		client: &http.Client{},
		secret: []byte("secret"),
		Logger: config.GetLogger().Logger,
	}
	tokenServiceData, err := sStub.GenerateToken(access, "user@example.com", 2048)
	assert.NoError(t, err)
	claims, err := sStub.CheckTokenValidness(tokenServiceData.Token)
	assert.NoError(t, err)
	flag, err := IsExpired(claims)
	assert.NoError(t, err)
	assert.False(t, flag)
}

func TestCheckTokenValidness_InvalidToken(t *testing.T) {
	invalidToken := "Z2xhZHlzLmNoYW1wbEBlZG1zLmNvbTphQG0xbg"
	sStub := SessionsService{
		client: &http.Client{},
		secret: []byte("secret"),
		Logger: config.GetLogger().Logger,
	}

	claims, err := sStub.CheckTokenValidness(invalidToken)
	assert.Error(t, err)
	flag, err := IsExpired(claims)
	assert.Error(t, err)
	assert.False(t, flag)
}
