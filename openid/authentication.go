package openid

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Authenticate(flow string, config Configuration) (string, error) {
	switch flow {
	case "code":
		at, err := code(config)
		return at.Access_token, err
	case "client":
		at, err := client_credentials(config)
		return at.Access_token, err
	default:
		return "", errors.New("unknown grant type")
	}
}

func code(config Configuration) (*AccessToken, error) {
	return nil, nil
}

func client_credentials(config Configuration) (*AccessToken, error) {
	fmt.Println("Starting client credentials request")
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	req, err := http.NewRequest("POST", config.OpenID.Token_endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	client_creds := []byte(config.ClientId + ":" + config.ClientSecret)
	encoded_client_creds := base64.RawURLEncoding.EncodeToString(client_creds)
	req.Header.Add("Authorization", "Basic "+encoded_client_creds)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var at AccessToken
	err = json.Unmarshal(b, &at)
	return &at, err
}
