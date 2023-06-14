package openid

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/safinwasi/loki/mocks"
	"github.com/stretchr/testify/assert"
)

func init() {
	Client = &mocks.MockHttpClient{}
}

func TestFetchOpenId(t *testing.T) {
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		openid_json := `{ "authorization_endpoint": "/authorize",
 "registration_endpoint": "/registration",
 "token_endpoint": "/token"}`
		body := io.NopCloser(bytes.NewReader([]byte(openid_json)))
		return &http.Response{StatusCode: 200, Body: body}, nil
	}
	_, err := Fetch_openid("https://account.google.com")
	assert.Equal(t, err, nil, "error should be nil")

}
