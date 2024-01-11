package openid

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func Register(hostname string, values RegistrationPayload) (*Configuration, error) {
	log.Println("Starting client registration request")
	_, err := url.Parse(hostname)
	if err != nil {
		return nil, err
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
