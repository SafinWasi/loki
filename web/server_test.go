package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/safinwasi/loki/secrets"
)

var testHost = ""

func TestHome(t *testing.T) {
	secrets.Initialize(false)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	homeHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestMissingClient(t *testing.T) {
	secrets.Initialize(false)
	req := httptest.NewRequest(http.MethodGet, "/client/test", nil)
	w := httptest.NewRecorder()
	clientHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestFoundClient(t *testing.T) {
	secrets.Initialize(false)
	secrets.Set("test", []byte("test"))
	req := httptest.NewRequest(http.MethodGet, "/client/test", nil)
	w := httptest.NewRecorder()
	clientHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestDeleteNotFound(t *testing.T) {
	secrets.Initialize(false)
	req := httptest.NewRequest(http.MethodGet, "/delete/test", nil)
	w := httptest.NewRecorder()
	deleteHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %d, got %d", http.StatusNotFound, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestDeleteFound(t *testing.T) {
	secrets.Initialize(false)
	req := httptest.NewRequest(http.MethodGet, "/delete/test", nil)
	w := httptest.NewRecorder()
	secrets.Set("test", []byte("test"))
	deleteHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Expected %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestRegistration(t *testing.T) {
	secrets.Initialize(false)
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
	data := url.Values{}
	data.Set("host", ts.URL)
	data.Set("payload", "")
	urlEncoded := data.Encode()
	reader := strings.NewReader(urlEncoded)
	req := httptest.NewRequest(http.MethodPost, "/delete/test", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	registrationHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestGetAfterCreation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	homeHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
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
