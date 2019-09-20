package db

import (
	"github.com/Soroka-EDMS/svc/sessions/pkgs/models"
	"github.com/go-kit/kit/log"
)

//Connection returns interface object that enclosing database after connection to it
func Connection(logger log.Logger, conn string) (models.ISessionDatabase, error) {
	var (
		db  SessionsDbStub
		err error
	)
	if conn == "stub" {
		db.Tokens = make(map[string]string)
		db.Logger = logger
	} else {
		//Real database
	}

	return &db, err
}
