package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/safinwasi/loki/openid"
)

func TestCodeUrlCreation(t *testing.T) {
	conf := openid.Configuration{
		Client_id:     "test",
		Client_secret: "test",
	}
	url := CreateCodeUrl(conf)
	expected := "client_id=test&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&response_type=code&scope=openid"
	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestTokenRequest(t *testing.T) {
	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/token":
			w.Write(getDummyToken())
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	token, err := SendTokenRequest("test", "test", "test", ts.URL+"/token", "code")
	if err != nil {
		t.Error(err)
	}
	if token != "test" {
		t.Errorf("Expected test, got %s", token)
	}
}

func getDummyToken() []byte {
	output :=
		`
	{
		"access_token": "test"
	}
`
	return []byte(output)
}
