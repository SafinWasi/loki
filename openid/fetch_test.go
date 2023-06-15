package openid

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/safinwasi/loki/mocks"
	"github.com/stretchr/testify/assert"
)

const hostName = "https://account.google.com"

func init() {
	Client = &mocks.MockHttpClient{}
}

func TestFetchOpenId(t *testing.T) {
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		openid_json := `{ "authorization_endpoint": "/authorize", "registration_endpoint": "/registration","token_endpoint": "/token"}`
		body := io.NopCloser(bytes.NewReader([]byte(openid_json)))
		return &http.Response{StatusCode: 200, Body: body}, nil
	}
	openid, err := Fetch_openid(hostName)
	assert.Equal(t, err, nil, "error should be nil")
	assert.Equal(t, openid.Hostname, hostName, "expected "+hostName)

}

func TestRegistration(t *testing.T) {
	ssa := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		openid_json := `{ "client_id": "abc", "client_secret": "def"}`
		body := io.NopCloser(bytes.NewReader([]byte(openid_json)))
		return &http.Response{StatusCode: 200, Body: body}, nil
	}
	config, err := Register(hostName, ssa)
	assert.Equal(t, err, nil, "error should be nil")
	assert.Equal(t, config.Client_id, "abc", "mismatch in client ID")
	assert.Equal(t, config.Client_secret, "def", "mismatch in client secret")
}

func TestClientCredentials(t *testing.T) {
	var config = Configuration{}
	var oidc = OIDCServer{}
	oidc.Token_endpoint = "/token"
	config.Client_id = "abcdef"
	config.Client_secret = "abcdef"
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		openid_json := `{ "access_token": "abc"}`
		body := io.NopCloser(bytes.NewReader([]byte(openid_json)))
		return &http.Response{StatusCode: 200, Body: body}, nil
	}
	token, err := Authenticate("client", config)
	assert.Equal(t, err, nil, "error should be nil")
	assert.Equal(t, token, "abc", "mismatch in token string (credentials)")
}
