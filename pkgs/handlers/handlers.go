package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

//MakeHTTPHandler wraps all service handlers in one HTTP handler
func MakeHTTPHandler(endp endpoints.SessionsEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET     /session/login                      generates a pair of tokens
	// GET     /session/logout                     restoke refresh token from cookie
	// PUT     /session/check_token                checks whether an access token is valid (contain valid priveledges and not expired). Regenerates token if expired

	r.Methods("GET").Path(constants.LoginEndpoint).Handler(httptransport.NewServer(
		endp.LoginEndpoint,
		DecodeLoginRequest,
		encodeLoginResponse,
		options...,
	))

	r.Methods("GET").Path(constants.LogoutEndpoint).Handler(httptransport.NewServer(
		endp.LogoutEndpoint,
		DecodeLogoutRequest,
		encodeLogoutResponse,
		options...,
	))

	r.Methods("POST").Path(constants.CheckTokenEndpoint).Handler(httptransport.NewServer(
		endp.CheckTokenEndpoint,
		DecodeCheckTokenRequest,
		encodeCheckTokenResponse,
		options...,
	))

	return r
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.LoginData
	user, pass, ok := r.BasicAuth()

	if !ok {
		return nil, errors.ErrNonAuthorized
	}

	req.UserName = user
	req.Password = pass

	return endpoints.LoginRequest{Req: req}, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(endpoints.LoginResponse)
	if !ok {
		return errors.ErrEncoding
	}

	err := e.Error()
	if err != nil {
		return err
	}

	AddCookie(w, e.RefreshToken.Token, e.RefreshToken.ExpirationDate)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(e.AccessToken)
}

func DecodeLogoutRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.LogoutData
	req.Cookie, err = r.Cookie("refresh_token")

	if err != nil {
		return nil, errors.ErrNonAuthorized
	}

	return endpoints.LogoutRequest{Req: req}, nil
}

func encodeLogoutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	AddCookie(w, "", 0)
	return nil
}

func DecodeCheckTokenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var reqAnotherService models.CheckTokenAnotherServiceInput
	var reqCurrentService models.CheckTokenServiceInput

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

	config.GetLogger().Logger.Log("a", reqCurrentService.AccessToken, "r", reqCurrentService.RefreshToken)

	return endpoints.CheckTokenRequest{Req: reqCurrentService}, nil
}

func encodeCheckTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(endpoints.CheckTokenResponse)
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

	json.NewEncoder(w).Encode(models.ErrorResponse{
		Reason:  errReason,
		Message: err.Error(),
	})
}

func codeFrom(err error) (int, string) {
	switch err {
	case errors.ErrMissingBody:
		return http.StatusBadRequest, constants.MissingBody
	case errors.ErrMalformedBody:
		return http.StatusBadRequest, constants.MalformedBody
	case errors.ErrMisingRefreshToken:
		return http.StatusUnauthorized, constants.MissingRefreshToken
	case errors.ErrInvalidClaimInToken:
		return http.StatusUnauthorized, constants.InvalidClaimInToken
	case errors.ErrNonAuthorized:
		return http.StatusUnauthorized, constants.NonAuthorized
	case errors.ErrEncoding:
		return http.StatusInternalServerError, constants.Encoding
	default:
		return http.StatusInternalServerError, err.Error()
	}
}

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
