package openid

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

func Request(disable_ssl bool, request *http.Request) ([]byte, error) {
	var client *http.Client
	if disable_ssl {
		customTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: customTransport}
	} else {
		client = &http.Client{}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	response_bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Request to %v, status %v", request.URL, response.Status)
	return response_bytes, err
}
