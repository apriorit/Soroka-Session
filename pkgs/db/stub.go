package db

import (
	"sync"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/errors"
)

type SessionsDbStub struct {
	Tokens map[string]string
	Mtx    sync.RWMutex
	Logger log.Logger
}

func (db *SessionsDbStub) Save(userID string, token string) (err error) {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()
	db.Tokens[userID] = token

	return nil
}

func (db *SessionsDbStub) Get(userID string) (token string, err error) {
	var ok bool

	db.Mtx.Lock()
	defer db.Mtx.Unlock()

	if token, ok = db.Tokens[userID]; !ok {
		return "", errors.ErrClientUnkown
	}

	return token, nil
}

func (db *SessionsDbStub) Exist(userID string, token string) (flag bool) {
	var (
		ok bool
		t  string
	)

	db.Mtx.Lock()
	defer db.Mtx.Unlock()

	if t, ok = db.Tokens[userID]; !ok || token != t {
		return false
	}

	return true
}

func (db *SessionsDbStub) Delete(userID string, token string) (err error) {
	db.Mtx.Lock()
	defer db.Mtx.Unlock()

	delete(db.Tokens, userID)
	return
}
