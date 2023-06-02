package openid

import (
	"encoding/json"
	"io"
	"net/http"
)

type OIDCServer struct {
	Authorization_endpoint string
	Token_endpoint         string
}

func Fetch_openid(hostname string) (*OIDCServer, error) {
	var oidc OIDCServer
	resp, err := http.Get(hostname + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &oidc)
	}
	return &oidc, nil
}
