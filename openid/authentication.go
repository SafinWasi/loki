package openid

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	// Credits to https://medium.com/@int128/shutdown-http-server-by-endpoint-in-go-2a0e2d7f9b8c
	authorization_uri := config.OpenID.Authorization_endpoint
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("client_id", config.Client_id)
	data.Set("redirect_uri", "http://localhost:8080/callback")
	data.Set("scope", "openid")
	fmt.Printf("Please visit %v\n", authorization_uri+"?"+data.Encode())
	m := http.NewServeMux()
	s := http.Server{Addr: ":8080", Handler: m}
	codeCh := make(chan string)
	m.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		// Send query parameter to the channel
		codeCh <- r.URL.Query().Get("code")
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	code := <-codeCh
	// Post process after shutdown here
	s.Shutdown(context.Background())
	log.Printf("Got code=%s", code)
	retrieve_token := func(code string) (*AccessToken, error) {
		token_data := url.Values{}
		token_data.Set("code", code)
		token_data.Set("grant_type", "authorization_code")
		token_data.Set("redirect_uri", "http://localhost:8080/callback")
		//token_data.Set("client_id", config.Client_id)
		req, err := http.NewRequest("POST", config.OpenID.Token_endpoint, bytes.NewBufferString(token_data.Encode()))
		if err != nil {
			return nil, err
		}
		client_creds := []byte(config.Client_id + ":" + config.Client_secret)
		encoded_client_creds := base64.RawURLEncoding.EncodeToString(client_creds)
		req.Header.Add("Authorization", "Basic "+encoded_client_creds)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		b, err := Request(disable_ssl, req)
		if err != nil {
			return nil, err
		}
		var at AccessToken
		err = json.Unmarshal(b, &at)
		return &at, err
	}
	return retrieve_token(code)
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
