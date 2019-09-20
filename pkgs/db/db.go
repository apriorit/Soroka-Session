package db

import (
	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	"github.com/go-kit/kit/log"
)

//Connection returns interface object that enclosing database after connection to it
func Connection(logger log.Logger, conn string) (models.ISessionDatabase, error) {
	if conn == "stub" {
		var db SessionsDbStub
		db.Tokens = make(map[string]string)
		db.Logger = logger
		return &db, nil
	} else {
		return nil, errors.ErrNotImplemented
	}
}
