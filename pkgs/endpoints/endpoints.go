package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//Endpoints collects individually constructed endpoints into a single type. Each endpoint is a func that wraps corresponding function from service interface
type SessionsEndpoints struct {
	LoginEndpoint      endpoint.Endpoint
	LogoutEndpoint     endpoint.Endpoint
	CheckTokenEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s models.ISessionService) SessionsEndpoints {
	return SessionsEndpoints{
		LoginEndpoint:      BuildLoginEndpoint(s),
		LogoutEndpoint:     BuildLogoutEndpoint(s),
		CheckTokenEndpoint: BuildCheckTokenEndpoint(s),
	}
}

func BuildLoginEndpoint(svc models.ISessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginRequest)
		at, rt, e := svc.Login(ctx, req.Req)
		return LoginResponse{AccessToken: at, RefreshToken: rt, Err: e}, nil
	}
}
func BuildLogoutEndpoint(svc models.ISessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LogoutRequest)
		e := svc.Logout(ctx, req.Req)
		return LogoutResponse{Err: e}, nil
	}
}
func BuildCheckTokenEndpoint(svc models.ISessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CheckTokenRequest)
		t, e := svc.CheckToken(ctx, req.Req)
		return CheckTokenResponse{Req: t, Err: e}, nil
	}
}

type LoginRequest struct {
	Req models.LoginData
}

type LogoutRequest struct {
	Req models.LogoutData
}

type LoginResponse struct {
	AccessToken  models.TokenData
	RefreshToken models.TokenData
	Err          error
}

type CheckTokenRequest struct {
	Req models.CheckTokenServiceInput
}

type CheckTokenResponse struct {
	Req models.CheckTokenServiceOutput
	Err error
}

type LogoutResponse struct {
	Err error
}

func (resp LoginResponse) Error() error      { return resp.Err }
func (resp LogoutResponse) Error() error     { return resp.Err }
func (resp CheckTokenResponse) Error() error { return resp.Err }
