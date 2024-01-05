package openid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testHost = ""

func TestWellknownRetrieval(t *testing.T) {
	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/.well-known/openid-configuration":
			w.Write(getDummyWellknown())
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	testHost = ts.URL
	_, err := Fetch_openid(ts.URL)
	if err != nil {
		t.Error(err)
	}
}

func TestRegistration(t *testing.T) {
	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/.well-known/openid-configuration":
			w.Write(getDummyWellknown())
		case "/register":
			w.Write(getDummyClient())
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	testHost = ts.URL
	_, err := Register(ts.URL, "")
	if err != nil {
		t.Error(err)
	}
	_, err = Register(ts.URL, getDummyPayload())
}

func getDummyWellknown() []byte {
	output := fmt.Sprintf(`
	{
		"hostname": "%s",
		"authorization_endpoint": "%s/authorize",
		"token_endpoint": "%s/token",
		"registration_endpoint": "%s/register"
	}
`, testHost, testHost, testHost, testHost)
	return []byte(output)
}

func getDummyClient() []byte {
	output := `
	{
		"client_id": "test",
		"client_secret": "test"
	}
`
	return []byte(output)
}

func getDummyPayload() string {
	output := `
	{
		"redirect_uris": ["localhost:3000"]
	}
`
	return output
}
