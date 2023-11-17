package openid

type Configuration struct {
	Client_id     string     `json:"client_id"`
	Client_secret string     `json:"client_secret"`
	Acr_values    []string   `json:"acr_values"`
	OpenID        OIDCServer `json:"openid"`
}

type OIDCServer struct {
	Hostname               string `json:"hostname"`
	Authorization_endpoint string `json:"authorization_endpoint"`
	Token_endpoint         string `json:"token_endpoint"`
	Registration_endpoint  string `json:"registration_endpoint"`
}

type AccessToken struct {
	Access_token string `json:"access_token"`
	scope        string
	token_type   string
	expires_in   string
}

type RegistrationPayload struct {
	Redirect_uris  []string `json:"redirect_uris"`
	Scope          []string `json:"scope"`
	Grant_types    []string `json:"grant_types"`
	Response_types []string `json:"response_types"`
	Client_name    string   `json:"client_name"`
	Ssa            string   `json:"software_statement"`
	Acr_values     []string `json:"acr_values"`
	Lifetime       int      `json:"lifetime"`
}
