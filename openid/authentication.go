package openid

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func Authenticate(flow string, config Configuration, disable_ssl bool) (string, error) {
	switch flow {
	case "code":
		at, err := code(config, disable_ssl)
		return at.Access_token, err
	case "client":
		at, err := client_credentials(config, disable_ssl)
		return at.Access_token, err
	default:
		return "", errors.New("unknown grant type")
	}
}

func code(config Configuration, disable_ssl bool) (*AccessToken, error) {
	return nil, nil
}

func client_credentials(config Configuration, disable_ssl bool) (*AccessToken, error) {
	fmt.Println("Starting client credentials request")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	req, err := http.NewRequest("POST", config.OpenID.Token_endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	client_creds := []byte(config.Client_id + ":" + config.Client_secret)
	encoded_client_creds := base64.RawURLEncoding.EncodeToString(client_creds)
	req.Header.Add("Authorization", "Basic "+encoded_client_creds)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := Request(disable_ssl, req)
	if err != nil {
		return nil, err
	}
	var at AccessToken
	err = json.Unmarshal(response, &at)
	return &at, err
}
