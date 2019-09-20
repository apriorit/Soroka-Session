package service

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

type TokenType int

const (
	access TokenType = iota
	refresh
)

type SessionsService struct {
	Db     models.ISessionDatabase
	client *http.Client
	secret []byte
	Logger log.Logger
}

//NewSessionsService creates session service
func NewSessionsService(db models.ISessionDatabase, s, key []byte) models.ISessionService {
	cl, _ := MakeHTTPClient(key)

	return &SessionsService{
		Db:     db,
		client: cl,
		secret: s,
		Logger: config.GetLogger().Logger,
	}
}

//Build creates session service with middleware
func Build(logger log.Logger, db models.ISessionDatabase, secret, pKey []byte) models.ISessionService {
	var svc models.ISessionService
	{
		svc = NewSessionsService(db, secret, pKey)
		svc = LoggingMiddleware(logger)(svc)
	}

	return svc
}

//Login handles login requets
func (svc *SessionsService) Login(ctx context.Context, ld models.LoginData) (resAccess, resRefresh models.TokenData, err error) {
	//1. Authentificate user
	err = svc.EnsureUserCreds(ld.UserName, ld.Password)
	if err != nil {
		svc.Logger.Log("method", "EnsureUserCreds", "action", "checking user credentials", "error", err)
		return resAccess, resRefresh, err
	}

	//2. Get user profile
	profile, err := svc.GetUserProfile(ld.UserName, string(svc.secret))
	if err != nil {
		svc.Logger.Log("method", "GetUserProfile", "action", "retrieving user profile", "error", err)
		return resAccess, resRefresh, err
	}

	//3. Create an access token
	resAccess, err = svc.GenerateToken(access, ld.UserName, profile.Role.Mask)
	if err != nil {
		svc.Logger.Log("method", "GetToken", "action", "generate access token", "error", err)
		return resAccess, resRefresh, err
	}

	resRefresh, err = svc.GenerateToken(refresh, ld.UserName, profile.Role.Mask)
	if err != nil {
		svc.Logger.Log("method", "GetToken", "action", "generate refresh token", "error", err)
		return resAccess, resRefresh, err
	}

	//4. Save token in db
	svc.Db.Save(ld.UserName, resRefresh.Token)

	return resAccess, resRefresh, nil
}

//Logout handles logout requets
func (svc *SessionsService) Logout(ctx context.Context, lod models.LogoutData) error {
	var (
		err           error
		ok            bool
		sub           string
		refreshClaims jwt.MapClaims
	)

	//If a refresh token is invalid then nothing to delete from database. Just logging out a user
	if refreshClaims, err = svc.CheckTokenValidness(lod.Cookie.Value); err != nil {
		svc.Logger.Log("method", "Logout", "err", errors.ErrClientUnknown)
		return err
	}

	if sub, ok = refreshClaims["sub"].(string); !ok {
		svc.Logger.Log("method", "Logout", "err", errors.ErrClientUnknown)
		return errors.ErrClientUnknown
	}
	svc.Db.Delete(sub, lod.Cookie.Value)
	return nil
}

//CheckToken checks whether an access token is valid and regenerates it if so
func (svc *SessionsService) CheckToken(ctx context.Context, td models.CheckTokenServiceInput) (res models.CheckTokenServiceOutput, err error) {
	//1. Check whether tokens are valid (signed with HMAC method and our service secret)
	accessTokenClaims, err := svc.CheckTokenValidness(td.AccessToken)
	if err != nil {
		return res, err
	}
	refreshTokenClaims, err := svc.CheckTokenValidness(td.RefreshToken)
	if err != nil {
		return res, err
	}

	//2. Check expiration claim
	atexp, err := IsExpired(accessTokenClaims)

	if err != nil {
		return res, err
	}
	rtexp, err := IsExpired(refreshTokenClaims)

	if err != nil {
		return res, err
	}

	if atexp && !rtexp {
		//Access token is expired, refresh token is not expired. Regenerate an access token if such token exist for current user
		var (
			ok   bool
			sub  string
			mask float64
		)
		if sub, ok = accessTokenClaims["sub"].(string); !ok {
			return res, errors.ErrInvalidClaimInToken
		}

		if mask, ok = accessTokenClaims["mask"].(float64); !ok {
			return res, errors.ErrInvalidClaimInToken
		}

		//Check refresh token existance
		if ok = svc.Db.Exist(sub, td.RefreshToken); !ok {
			return res, errors.ErrNonAuthorized
		}

		tokenData, err := svc.GenerateToken(access, sub, int64(mask))
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
		err = errors.ErrNonAuthorized
	}

	return res, err
}
