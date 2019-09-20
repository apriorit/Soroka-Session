package errors

import (
	"errors"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
)

var (
	ErrMissingBody          = errors.New(constants.MissingBody)
	ErrMisingRefreshToken   = errors.New(constants.MissingRefreshToken)
	ErrExpiredRefreshToken  = errors.New(constants.ExpiredRefreshToken)
	ErrExpiredAccessToken   = errors.New(constants.ExpiredAccessToken)
	ErrNoPermissions        = errors.New(constants.NoPermissions)
	ErrMalformedBody        = errors.New(constants.MalformedBody)
	ErrEncoding             = errors.New(constants.Encoding)
	ErrNonAuthorized        = errors.New(constants.NonAuthorized)
	ErrRequestToUsersFailed = errors.New(constants.RequestToUsersFailed)
	ErrClientUnkown         = errors.New(constants.ClientUnkown)
	ErrInvalidClaimInToken  = errors.New(constants.InvalidClaimInToken)
	ErrFailedToCreateJWT    = errors.New(constants.FailedToCreateJWT)
	ErrPublicKeyIsMissing   = errors.New(constants.PublicKeyIsMissing)
	ErrInvalidTokenType     = errors.New(constants.InvalidTokenType)
)
