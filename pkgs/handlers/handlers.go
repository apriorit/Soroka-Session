package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	conf "github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	c "github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
	e "github.com/Soroka-EDMS/svc/sessions/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//MakeHTTPHandler wraps all service handlers in one HTTP handler
func MakeHTTPHandler(endp e.SessionsEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET     /session/login                      generates a pair of tokens
	// GET     /session/logout                     restoke refresh token from cookie
	// PUT     /session/check_token                checks whether an access token is valid (contain valid priveledges and not expired). Regenerates token if expired

	r.Methods("GET").Path(c.LoginEndpoint).Handler(httptransport.NewServer(
		endp.LoginEndpoint,
		DecodeLoginRequest,
		encodeLoginResponse,
		options...,
	))

	r.Methods("GET").Path(c.LogoutEndpoint).Handler(httptransport.NewServer(
		endp.LogoutEndpoint,
		DecodeLogoutRequest,
		encodeLogoutResponse,
		options...,
	))

	r.Methods("POST").Path(c.CheckTokenEndpoint).Handler(httptransport.NewServer(
		endp.CheckTokenEndpoint,
		DecodeCheckTokenRequest,
		encodeCheckTokenResponse,
		options...,
	))

	return r
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req m.LoginData
	user, pass, ok := r.BasicAuth()

	if !ok {
		return nil, errors.ErrNonAuthorized
	}

	req.UserName = user
	req.Password = pass

	return e.LoginRequest{Req: req}, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(e.LoginResponse)
	if !ok {
		return errors.ErrEncoding
	}

	err := e.Error()
	if err != nil {
		return err
	}

	AddCookie(w, "new", &e.RefreshToken)
	conf.GetLogger().Logger.Log("refresh_token", e.RefreshToken.Token)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.AccessToken)
}

func DecodeLogoutRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req m.LogoutData
	req.Cookie, err = r.Cookie("refresh_token")

	if err != nil {
		return nil, errors.ErrNonAuthorized
	}

	return e.LogoutRequest{Req: req}, nil
}

func encodeLogoutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	AddCookie(w, "expired", nil)
	return nil
}

func DecodeCheckTokenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqAnotherService m.CheckTokenAnotherServiceInput
	var reqCurrentService m.CheckTokenServiceInput

	//Parse request body and get access token
	if r.Body == nil {
		return nil, errors.ErrMissingBody
	}
	if err := json.NewDecoder(r.Body).Decode(&reqAnotherService); err != nil {
		return nil, err
	}
	reqCurrentService.AccessToken = reqAnotherService.AccessToken

	//Parse cookie and get refresh token
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, errors.ErrMisingRefreshToken
	}
	reqCurrentService.RefreshToken = cookie.Value

	conf.GetLogger().Logger.Log("a", reqCurrentService.AccessToken, "r", reqCurrentService.RefreshToken)

	return e.CheckTokenRequest{Req: reqCurrentService}, nil
}

func encodeCheckTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(e.CheckTokenResponse)
	if !ok {
		return errors.ErrEncoding
	}

	err := e.Error()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.Req)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	errCode, errReason := codeFrom(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errCode)

	json.NewEncoder(w).Encode(m.ErrorResponse{
		Reason:  errReason,
		Message: err.Error(),
	})
}

func codeFrom(err error) (int, string) {
	switch err {
	case errors.ErrMissingBody:
		return http.StatusBadRequest, c.MissingBody
	case errors.ErrMalformedBody:
		return http.StatusBadRequest, c.MalformedBody
	case errors.ErrMisingRefreshToken:
		return http.StatusUnauthorized, c.MissingRefreshToken
	case errors.ErrNonAuthorized:
		return http.StatusUnauthorized, c.NonAuthorized
	case errors.ErrEncoding:
		return http.StatusInternalServerError, c.Encoding
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
