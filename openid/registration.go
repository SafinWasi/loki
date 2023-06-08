package openid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Register(hostname string, ssa string) (*Configuration, error) {
	fmt.Println("Starting client registration request")
	values := make(map[string]any)
	values["redirect_uris"] = []string{"http://localhost:8080/callback"}
	values["scope"] = []string{"openid", "profile"}
	values["grant_types"] = []string{"authorization_code", "client_credentials"}
	values["response_types"] = []string{"code", "token"}
	values["client_name"] = "loki_client"
	if ssa != "" {
		values["software_statement"] = ssa
		values["redirect_uris"] = []string{hostname}
	}
	body_bytes, err := json.MarshalIndent(values, "", "\t")
	if err != nil {
		return nil, err
	}
	oidc, err := Fetch_openid(hostname)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", oidc.Registration_endpoint, bytes.NewReader(body_bytes))
	if err != nil {
		return nil, err
	}
	response, err := Request(request)
	if err != nil {
		return nil, err
	}
	var new_client Configuration
	json.Unmarshal(response, &new_client)
	new_client.OpenID = *oidc
	return &new_client, nil
}
