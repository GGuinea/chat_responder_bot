package auth

const (
	GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"
)

type accessTokenGrantRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func NewDefaultAccessTokenGrantRequest(code, clientId, clientSecret, redirectUrl string) accessTokenGrantRequest {
	return accessTokenGrantRequest{
		GrantType:    GRANT_TYPE_AUTHORIZATION_CODE,
		Code:         code,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURI:  redirectUrl,
	}
}

type AccessTokenGrantResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
