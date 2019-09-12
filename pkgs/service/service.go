package service

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"

	cfg "github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	e "github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

type SessionsServiceStub struct {
	client *http.Client
	secret string
	Logger log.Logger
}

//NewSessionsService creates session service
func NewSessionsService(s string) m.SessionService {
	return &SessionsServiceStub{
		client: &http.Client{},
		secret: s,
		Logger: cfg.GetLogger().Logger,
	}
}

//Build creates session service with middleware
func Build(logger log.Logger, s string) m.SessionService {
	var svc m.SessionService
	{
		svc = NewSessionsService(s)
		svc = LoggingMiddleware(logger)(svc)
	}

	return svc
}

//Login handles login requets
func (sStub *SessionsServiceStub) Login(ctx context.Context, ld m.LoginData) (resAccess, resRefresh m.TokenData, err error) {
	//1. Authentificate user
	err = sStub.EnsureUserCreds(ld.UserName, ld.Password)
	if err != nil {
		sStub.Logger.Log("method", "EnsureUserCreds", "action", "checking user credentials", "error", err)
		return resAccess, resRefresh, err
	}

	//2. Get user profile
	profile, err := sStub.GetUserProfile(ld.UserName, sStub.secret)
	if err != nil {
		sStub.Logger.Log("method", "GetUserProfile", "action", "retrieving user profile", "error", err)
		return resAccess, resRefresh, err
	}

	//3. Create an access token
	resAccess, err = sStub.GenerateToken("access", ld.UserName, profile.Role.Mask)
	if err != nil {
		sStub.Logger.Log("method", "GetToken", "action", "generate access token", "error", err)
		return resAccess, resRefresh, err
	}

	resRefresh, err = sStub.GenerateToken("refresh", ld.UserName, profile.Role.Mask)
	if err != nil {
		sStub.Logger.Log("method", "GetToken", "action", "generate refresh token", "error", err)
		return resAccess, resRefresh, err
	}

	return resAccess, resRefresh, nil
}

//Logout handles logout requets
func (sStub *SessionsServiceStub) Logout(ctx context.Context, lod m.LogoutData) error {
	//1. Get refresh token value from cookie
	//2. Parse refresh token
	//3. Restore session sensetive data from session database
	//4. In handler layer decodeLogourResponse will set cookie within invalid refresh token
	sStub.Logger.Log("Method", "Logout", "message", "service handler has been reached")
	return nil
}

//CheckToken checks whether an access token is valid and regenerates it if so
func (sStub *SessionsServiceStub) CheckToken(ctx context.Context, td m.CheckTokenServiceInput) (res m.CheckTokenServiceOutput, err error) {
	//1. Check whether tokens are valid (signed with HMAC method and our service secret)
	accessTokenClaims, err := sStub.CheckTokenValidness(td.AccessToken)
	if err != nil {
		return res, err
	}

	refreshTokenClaims, err := sStub.CheckTokenValidness(td.RefreshToken)
	if err != nil {
		return res, err
	}

	//2. Check expiration claim
	atexp := IsExpired(accessTokenClaims)
	rtexp := IsExpired(refreshTokenClaims)
	if atexp && !rtexp {
		//Access token is expired, refresh token is not expired. Regenerate an access token
		sub := accessTokenClaims["sub"].(string)
		mask := accessTokenClaims["mask"].(float64)
		tokenData, err := sStub.GenerateToken("access", sub, int64(mask))
		if err != nil {
			return res, err
		}
		res.AccessToken = tokenData.Token //a new token
	} else if !atexp && !rtexp {
		//Access token is not expired, refresh token is not expired. Return an old access token
		res.AccessToken = td.AccessToken
	} else {
		//Refresh token is expired
		res.AccessToken = ""
		err = e.ErrNonAuthorized
	}

	return res, err
}
