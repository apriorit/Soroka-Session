package models

type ISessionDatabase interface {
	Save(userID string, token string) (err error)
	Exist(userID string, token string) (flag bool)
	Get(userID string) (token string, err error)
	Delete(userID string, token string) (err error)
}
