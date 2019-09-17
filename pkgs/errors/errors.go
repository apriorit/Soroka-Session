package errors

import (
	"errors"

	c "github.com/Soroka-EDMS/svc/sessions/pkgs/constants"
)

var (
	ErrMissingBody          = errors.New(c.MissingBody)
	ErrMisingRefreshToken   = errors.New(c.MissingRefreshToken)
	ErrExpiredRefreshToken  = errors.New(c.ExpiredRefreshToken)
	ErrExpiredAccessToken   = errors.New(c.ExpiredAccessToken)
	ErrNoPermissions        = errors.New(c.NoPermissions)
	ErrMalformedBody        = errors.New(c.MalformedBody)
	ErrEncoding             = errors.New(c.Encoding)
	ErrNonAuthorized        = errors.New(c.NonAuthorized)
	ErrRequestToUsersFailed = errors.New(c.RequestToUsersFailed)
	ErrClientUnkown         = errors.New(c.ClientUnkown)
	ErrFailedToCreateJWT    = errors.New(c.FailedToCreateJWT)
	ErrPublicKeyIsMissing   = errors.New(c.PublicKeyIsMissing)
)
