package openid

import (
	"errors"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

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
	if response.StatusCode >= http.StatusBadRequest {
		log.Println("Error in request")
		return nil, errors.New(response.Status)
	}
	return response_bytes, err
}

func init() {
	Client = &http.Client{}
}
