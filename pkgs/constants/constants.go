package constants

const (
	LoginEndpoint         string = "/session/login"
	LogoutEndpoint        string = "/session/logout"
	CheckTokenEndpoint    string = "/session/check_token"
	URIOnGetProfile       string = "https://users_users.service_1:443/user?email=%v"
	URIOnAuthentification string = "https://users_users.service_1:443/users/check_auth"
	TokenIssuer           string = "https://edms.com/sessions"
	PublicKeyIsMissing    string = "Public key is missing"
	RequestToUsersFailed  string = "Request to Users service failed"
	MissingBody           string = "Missing content in request body"
	MissingRefreshToken   string = "Missing refresh token"
	ExpiredRefreshToken   string = "Refresh token expired"
	ExpiredAccessToken    string = "Access token expired"
	NoPermissions         string = "Client does not have required permissions"
	MalformedBody         string = "Malformed content in request body"
	Encoding              string = "An error occured while enconding response"
	NonAuthorized         string = "Required authorization"
	ClientUnkown          string = "Client is unknown"
	FailedToCreateJWT     string = "Failed to create JWT"
	InvalidClaimInToken   string = "Invalid claim in token"
	InvalidTokenType      string = "Invalid token type"
	ClaimsPerAccessToken  int    = 7
	ClaimsPerRefreshToken int    = 6
)
