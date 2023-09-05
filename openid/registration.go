package openid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Register(hostname string, ssa string, payload_file string) (*Configuration, error) {
	fmt.Println("Starting client registration request")
	var values = RegistrationPayload{}
	if payload_file == "" {
		values.Redirect_uris = []string{"http://localhost:8080/callback"}
		values.Scope = []string{"openid", "profile"}
		values.Grant_types = []string{"authorization_code", "client_credentials"}
		values.Response_types = []string{"code", "token"}
		values.Client_name = "loki_client"
		if ssa != "" {
			values.Ssa = ssa
			values.Redirect_uris = []string{hostname}
		}
	} else {
		b, err := os.ReadFile(payload_file)
		log.Println(b)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &values)
		log.Println(values)
		if err != nil {
			return nil, err
		}
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
