package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	b "encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"

	c "github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
	e "github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//EnsureUserCreds sends authentification request to Users service
func (sStub *SessionsServiceStub) EnsureUserCreds(username, password string) (err error) {

	req, err := http.NewRequest("GET", c.URIOnAuthentification, nil)
	req.SetBasicAuth(username, password)

	if err != nil {
		return err
	}

	resp, err := sStub.client.Do(req)
	if err != nil {
		return e.ErrRequestToUsersFailed
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return e.ErrNonAuthorized
	case http.StatusNotFound:
		return e.ErrClientUnkown
	default:
		return fmt.Errorf("Authentification request failed with code: %v", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("Authentification request received content type: %v", resp.Header.Get("Content-Type"))
	}

	return nil
}

//GetUserProfile query user profile by email from Users database
func (sStub *SessionsServiceStub) GetUserProfile(email, token string) (profile m.UserProfile, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(c.URIOnGetProfile, email), nil)
	if err != nil {
		return m.UserProfile{}, err
	}

	req.Header.Add("Bearer", token)
	resp, err := sStub.client.Do(req)
	if err != nil {
		return m.UserProfile{}, e.ErrRequestToUsersFailed
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return m.UserProfile{}, e.ErrRequestToUsersFailed
	case http.StatusUnauthorized:
		return m.UserProfile{}, e.ErrNonAuthorized
	case http.StatusNotFound:
		return m.UserProfile{}, e.ErrClientUnkown
	default:
		return m.UserProfile{}, fmt.Errorf("Get profile request failed with code: %v", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return m.UserProfile{}, fmt.Errorf("Get profile received content type: %v", resp.Header.Get("Content-Type"))
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&profile)

	if err != nil {
		return m.UserProfile{}, fmt.Errorf("Profile decoding failed")
	}

	return profile, nil
}

//GenerateToken generates and signs token according to toke type. Uses sgining method based on HS256
func (sStub *SessionsServiceStub) GenerateToken(tokenType, id string, mask int64) (m.TokenData, error) {
	var err error

	sStub.Logger.Log("Method", "GenerateToken", "Sign secret", sStub.secret)

	claims, exp := CreatePayload(tokenType, id, c.TokenIssuer, mask)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sStub.secret))

	if err != nil {
		return m.TokenData{}, err
	}

	return m.TokenData{
		Token:          tokenString,
		Type:           "Bearer",
		ExpirationDate: exp,
	}, err
}

//CreatePayload returns claims for JWT. If token type is "access" it returns claims for access token, or for refresh token otherwise
func CreatePayload(tokenType, cid, iss string, mask int64) (jwt.MapClaims, int64) {
	var (
		exp int64
		iat int64
	)

	iat = time.Now().Unix()
	if tokenType == "access" {
		exp = time.Now().Add(time.Duration(24) * time.Hour).Unix()
		return jwt.MapClaims{
			"iss":  iss,                //token issue
			"sub":  cid,                //user email
			"iat":  iat,                //issued at
			"nbf":  iat,                //issued not before
			"exp":  exp,                //expiration time
			"mask": mask,               //user mask
			"aud":  []string{cid, iss}, //audience claim. See: https://tools.ietf.org/html/rfc7519#
		}, exp
	}

	exp = time.Now().Add(time.Duration(720) * time.Hour).Unix()
	return jwt.MapClaims{
		"iss": iss,
		"sub": cid,
		"iat": iat,
		"nbf": iat,
		"exp": exp,
		"aud": []string{cid, iss},
	}, exp
}

//CheckTokenValidness parses token and return its claims if token is valid. There is a 'Valid' flag in token that is triggered when using jwt.Parse
func (sStub *SessionsServiceStub) CheckTokenValidness(tokenString string) (jwt.MapClaims, error) {
	//Prepare key func
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//return session secret that was used in signing process
		return []byte(sStub.secret), nil
	}

	//Get raw token
	token, err := jwt.Parse(tokenString, keyFunc)

	//Check token validness
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	//Log and return error of token is invalid
	sStub.Logger.Log("method", "GetRawToken", "err", err)
	return jwt.MapClaims{}, e.ErrNonAuthorized
}

//IsExpired checks whether a token is expired
func IsExpired(claims jwt.MapClaims) bool {
	return claims.Valid() != nil
}

//EncodeSessionSecret encodes session secret to base64 string
func EncodeSessionSecret(s string) (encoded string) {
	return b.StdEncoding.EncodeToString([]byte(s))
}
