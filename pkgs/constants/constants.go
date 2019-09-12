package constants

const (
	LoginEndpoint         string = "/session/login"
	LogoutEndpoint        string = "/session/logout"
	CheckTokenEndpoint    string = "/session/check_token"
	URIOnGetProfile       string = "https://edms.com/api/v1/user?email=%v"
	URIOnAuthentification string = "https://edms.com/api/v1/user_auth"
	TokenIssuer           string = "https://edms.com/sessions"

	RequestToUsersFailed  string = "Request to Users service failed"
	MissingBody           string = "Missing content in request body"
	MissingRefreshToken   string = "Missing refresh token"
	ExpiredRefreshToken   string = "Refresh token expired"
	ExpiredAccessToken    string = "Access token expired"
	NoPermissions         string = "Client does not have required permissions"
	MalformedBody         string = "Malformed content in request body"
	Encoding              string = "An error occured while enconding response"
	NonAuthorized         string = "Required authorization"
	ClientUnkown          string = "Client are unknown"
	FailedToCreateJWT     string = "Failed to create JWT"
	ClaimsPerAccessToken  int    = 7
	ClaimsPerRefreshToken int    = 6
)
