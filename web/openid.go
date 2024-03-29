package web

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/safinwasi/loki/openid"
)

func CreateCodeUrl(config openid.Configuration, params string, acr string) string {
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("client_id", config.Client_id)
	data.Set("redirect_uri", "http://127.0.0.1:3000/callback")
	data.Set("scope", "openid")
	var paramMap map[string]string
	err := json.Unmarshal([]byte(params), &paramMap)
	if err != nil {
		return "Error unmarshaling JSON params. Please ensure valid JSON is provided."
	}
	for key, val := range paramMap {
		data.Set(key, val)
	}
	if acr != "" {
		data.Set("acr_values", acr)
	}
	uri := data.Encode()
	return uri
}

func SendTokenRequest(code string, client_id string, client_secret string, token_endpoint string, grant_type string, scope string) (string, error) {

	data := url.Values{}
	data.Set("scope", scope)
	data.Set("grant_type", grant_type)
	if grant_type == "authorization_code" {
		data.Set("code", code)
		data.Set("redirect_uri", "http://localhost:3000/callback")
	}
	req, err := http.NewRequest("POST", token_endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	client_creds := []byte(client_id + ":" + client_secret)
	encoded_client_creds := base64.RawURLEncoding.EncodeToString(client_creds)
	req.Header.Add("Authorization", "Basic "+encoded_client_creds)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := openid.Request(req)
	if err != nil {
		return "", err
	}
	var at openid.AccessToken
	err = json.Unmarshal(response, &at)
	if err != nil {
		return "", err
	}
	return at.Access_token, nil
}
