package openid

import (
	"io"
	"log"
	"net/http"
)

var Client *http.Client

func Request(request *http.Request) ([]byte, error) {
	response, err := Client.Do(request)
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

func init() {
	Client = &http.Client{}
}
