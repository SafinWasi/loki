package openid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

//func Register(hostname string, ssa string, payload_string string) (*Configuration, error) {
//	log.Println("Starting client registration request")
//	_, err := url.Parse(hostname)
//	if err != nil {
//		return nil, err
//	}
//	var values = RegistrationPayload{}
//	if payload_string == "" {
//		values.Redirect_uris = []string{"http://localhost:3000/callback"}
//		values.Scope = []string{"openid", "profile"}
//		values.Grant_types = []string{"authorization_code", "client_credentials"}
//		values.Response_types = []string{"code", "token"}
//		values.Client_name = "loki_client"
//		values.Lifetime = 86400
//	} else {
//		b := []byte(payload_string)
//		err = json.Unmarshal(b, &values)
//		if err != nil {
//			return nil, err
//		}
//	}
//	if len(ssa) > 0 {
//		values.Ssa = ssa
//		values.Redirect_uris = []string{hostname}
//	}
//	body_bytes, err := json.MarshalIndent(values, "", "\t")
//	if err != nil {
//		return nil, err
//	}
//	oidc, err := Fetch_openid(hostname)
//	if err != nil {
//		return nil, err
//	}
//	request, err := http.NewRequest("POST", oidc.Registration_endpoint, bytes.NewReader(body_bytes))
//	if err != nil {
//		return nil, err
//	}
//	response, err := Request(request)
//	if err != nil {
//		return nil, err
//	}
//	var new_client Configuration
//	json.Unmarshal(response, &new_client)
//	new_client.OpenID = *oidc
//	return &new_client, nil
//}

func Register(hostname string, payload string) (*Configuration, error) {
	log.Println("Starting client registration request")
	_, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}
	var values = RegistrationPayload{}
	if payload == "" {
		values.Redirect_uris = []string{"http://localhost:3000/callback"}
		values.Scope = []string{"openid", "profile"}
		values.Grant_types = []string{"authorization_code", "client_credentials"}
		values.Response_types = []string{"code", "token"}
		values.Client_name = fmt.Sprintf("loki_client_%d", time.Now().Unix())
		values.Lifetime = 86400
	} else {
		b := []byte(payload)
		err = json.Unmarshal(b, &values)
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
