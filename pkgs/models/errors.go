package models

type ErrorResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type MissingRefresh struct {
	AccessToken string `json:"access_token"`
}
