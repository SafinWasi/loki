package openid

type Configuration struct {
	Client_id     string
	Client_secret string
	Acr_values    []string
	OpenID        OIDCServer
}

type OIDCServer struct {
	Hostname               string
	Authorization_endpoint string
	Token_endpoint         string
	Registration_endpoint  string
}

type AccessToken struct {
	Access_token string
	scope        string
	token_type   string
	expires_in   string
}

type RegistrationPayload struct {
	redirect_uris  []string
	scope          []string
	grant_types    []string
	response_types []string
	client_name    string
	ssa            string
	acr_values     []string
}
