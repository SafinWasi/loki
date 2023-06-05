package openid

type Configuration struct {
	Hostname     string
	ClientId     string
	ClientSecret string
	OpenID       OIDCServer
}

type OIDCServer struct {
	Authorization_endpoint string
	Token_endpoint         string
}

type AccessToken struct {
	Access_token string
	scope        string
	token_type   string
	expires_in   string
}
