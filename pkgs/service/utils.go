package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//EnsureUserCreds sends authentification request to Users service
func (sStub *SessionsService) EnsureUserCreds(username, password string) (err error) {
	req, err := http.NewRequest("GET", constants.URIOnAuthentification, nil)
	req.SetBasicAuth(username, password)

	if err != nil {
		return err
	}

	resp, err := sStub.client.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return errors.ErrNonAuthorized
	case http.StatusNotFound:
		return errors.ErrClientUnkown
	default:
		return fmt.Errorf("Authentification request failed with code: %v", resp.StatusCode)
	}

	return nil
}

//GetUserProfile sends a request to Users service to obtain user profile by his email
func (sStub *SessionsService) GetUserProfile(email, token string) (profile models.UserProfile, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(constants.URIOnGetProfile, email), nil)
	if err != nil {
		return models.UserProfile{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := sStub.client.Do(req)
	if err != nil {
		return models.UserProfile{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		return models.UserProfile{}, errors.ErrRequestToUsersFailed
	case http.StatusUnauthorized:
		return models.UserProfile{}, errors.ErrNonAuthorized
	case http.StatusNotFound:
		return models.UserProfile{}, errors.ErrClientUnkown
	default:
		return models.UserProfile{}, fmt.Errorf("Get profile request failed with code: %v", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return models.UserProfile{}, fmt.Errorf("Get profile received content type: %v", resp.Header.Get("Content-Type"))
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&profile)

	if err != nil {
		return models.UserProfile{}, fmt.Errorf("Profile decoding failed")
	}

	return profile, nil
}

//GenerateToken generates and signs token according to toke type. Uses sgining method based on HS256
func (sStub *SessionsService) GenerateToken(tokenType TokenType, id string, mask int64) (models.TokenData, error) {
	var err error

	claims, exp, err := CreatePayload(tokenType, id, constants.TokenIssuer, mask)
	if err != nil {
		return models.TokenData{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sStub.secret))

	if err != nil {
		return models.TokenData{}, err
	}

	return models.TokenData{
		Token:          tokenString,
		Type:           "Bearer",
		ExpirationDate: exp,
	}, err
}

//CreatePayload returns claims for JWT. If token type is "access" it returns claims for access token, or for refresh token otherwise
func CreatePayload(tokenType TokenType, cid, iss string, mask int64) (jwt.MapClaims, int64, error) {
	var (
		exp    int64
		iat    int64
		claims jwt.MapClaims
	)

	iat = time.Now().Unix()

	switch tokenType {
	case access:
		exp = time.Now().Add(time.Duration(24) * time.Hour).Unix()
		claims = jwt.MapClaims{
			"iss":  iss,                //token issue
			"sub":  cid,                //user email
			"iat":  iat,                //issued at
			"nbf":  iat,                //issued not before
			"exp":  exp,                //expiration time
			"mask": mask,               //user mask
			"aud":  []string{cid, iss}, //audience claim. See: https://tools.ietf.org/html/rfc7519#
		}
	case refresh:
		exp = time.Now().Add(time.Duration(720) * time.Hour).Unix()
		claims = jwt.MapClaims{
			"iss": iss,                //token issue
			"sub": cid,                //user email
			"iat": iat,                //issued at
			"nbf": iat,                //issued not before
			"exp": exp,                //expiration time
			"aud": []string{cid, iss}, //audience claim. See: https://tools.ietf.org/html/rfc7519#
		}
	default:
		return jwt.MapClaims{}, 0, errors.ErrInvalidTokenType
	}

	return claims, exp, nil
}

//CheckTokenValidness parses token and return its claims if token is valid. There is a 'Valid' flag in token that is triggered when using jwt.Parse
func (sStub *SessionsService) CheckTokenValidness(tokenString string) (jwt.MapClaims, error) {
	//Get raw token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//return session secret that was used in signing process
		return []byte(sStub.secret), nil
	})

	if token == nil {
		return jwt.MapClaims{}, errors.ErrNonAuthorized
	}

	//Check token validness
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	//Log and return error if token is invalid
	sStub.Logger.Log("method", "GetRawToken", "err", err)
	return jwt.MapClaims{}, errors.ErrNonAuthorized
}

//IsExpired checks whether a token is expired
func IsExpired(claims jwt.MapClaims) (bool, error) {
	if claims == nil {
		return true, errors.ErrInvalidClaimInToken
	}

	exp := claims["exp"]
	var expValue int64

	switch v := exp.(type) {
	case float64:
		expValue = int64(v)
	case int64:
		expValue = v
	default:
		return false, errors.ErrInvalidClaimInToken
	}

	return time.Now().Unix() > expValue, nil
}

//EncodeSessionSecret encodes session secret to base64 string
func EncodeSessionSecret(s string) (encoded string) {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func MakeHTTPClient(pKey []byte) (*http.Client, error) {
	if len(pKey) == 0 {
		return nil, errors.ErrPublicKeyIsMissing
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(pKey)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true,
			},
		},
	}

	return client, nil
}
