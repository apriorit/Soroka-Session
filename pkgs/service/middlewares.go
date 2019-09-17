package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	m "github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(m.SessionService) m.SessionService

type loggingMiddleware struct {
	next   m.SessionService
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next m.SessionService) m.SessionService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (lmw loggingMiddleware) Login(ctx context.Context, ld m.LoginData) (resA m.TokenData, resR m.TokenData, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "Login", "usename", ld.UserName, "password", ld.Password, "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.Login(ctx, ld)
}

func (lmw loggingMiddleware) Logout(ctx context.Context, lod m.LogoutData) (err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "Logout", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.Logout(ctx, lod)
}

func (lmw loggingMiddleware) CheckToken(ctx context.Context, ctd m.CheckTokenServiceInput) (res m.CheckTokenServiceOutput, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "CheckToken", "token", ctd.AccessToken, "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.CheckToken(ctx, ctd)
}
