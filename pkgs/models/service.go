package models

//  service.go
//  https://github.com/apriorit/Soroka-EDMS/svc/users/pkgs/models
//
//  Created by Ivan Kashuba on 2019.09.03
//  Describe service models
import (
	"context"
	"net/http"
)

type ISessionService interface {
	Login(cntx context.Context, request LoginData) (resAccess, resRefresh TokenData, err error)
	Logout(cntx context.Context, request LogoutData) error
	CheckToken(cntx context.Context, request CheckTokenServiceInput) (res CheckTokenServiceOutput, err error)
}

type LoginData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LogoutData struct {
	Cookie *http.Cookie
}

type CheckTokenAnotherServiceInput struct {
	AccessToken string `json:"access_token"`
}

type CheckTokenServiceInput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CheckTokenServiceOutput struct {
	AccessToken string `json:"access_token"`
}

type TokenData struct {
	Token          string `json:"access_token"`
	Type           string `json:"type"`
	ExpirationDate int64  `json:"expiration_date"`
}

type UserRole struct {
	Name string
	Mask int64
}

type UserProfile struct {
	First_name    string   `json:"first_name"`
	Last_name     string   `json:"last_name"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	Location      string   `json:"location"`
	Position      string   `json:"position"`
	Status        bool     `json:"status"`
	Creation_date int64    `json:"creation_date"`
	Role          UserRole `json:"role"`
}
