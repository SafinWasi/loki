package openid

import (
	"encoding/json"
	"net/http"
)

func Fetch_openid(hostname string) (*OIDCServer, error) {
	var oidc OIDCServer
	request, err := http.NewRequest("GET", hostname+"/.well-known/openid-configuration", nil)
	if err != nil {
		return nil, err
	}
	resp, err := Request(request)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(resp, &oidc)
		if err != nil {
			return nil, err
		}
	}
	oidc.Hostname = hostname
	return &oidc, nil
}
