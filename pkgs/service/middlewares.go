package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(models.ISessionService) models.ISessionService

type loggingMiddleware struct {
	next   models.ISessionService
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next models.ISessionService) models.ISessionService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (lmw loggingMiddleware) Login(ctx context.Context, ld models.LoginData) (resA models.TokenData, resR models.TokenData, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "Login", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.Login(ctx, ld)
}

func (lmw loggingMiddleware) Logout(ctx context.Context, lod models.LogoutData) (err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "Logout", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.Logout(ctx, lod)
}

func (lmw loggingMiddleware) CheckToken(ctx context.Context, ctd models.CheckTokenServiceInput) (res models.CheckTokenServiceOutput, err error) {
	defer func(begin time.Time) {
		lmw.logger.Log("method", "CheckToken", "took", time.Since(begin), "err", err)
	}(time.Now())
	return lmw.next.CheckToken(ctx, ctd)
}
