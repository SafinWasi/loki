package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	req := httptest.NewRequest(http.MethodDelete, "/delete/test", nil)
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
	req := httptest.NewRequest(http.MethodDelete, "/delete/test", nil)
	w := httptest.NewRecorder()
	secrets.Set("test", []byte("test"))
	deleteHandler()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
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
		case "/token":
			w.Write(getDummyToken())
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	testHost = ts.URL
	data := url.Values{}
	data.Set("host", ts.URL)
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

func TestCodeFlow(t *testing.T) {
	keys, err := secrets.GetKeys()
	if err != nil {
		t.Error(err)
	}
	key := keys[0]
	req := httptest.NewRequest(http.MethodGet, "/code/"+key, nil)
	w := httptest.NewRecorder()
	codeFlow()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusFound {
		t.Errorf("Expected %d, got %d", http.StatusFound, res.StatusCode)
	}
	t.Cleanup(func() {
		secrets.RemoveKeyring()
	})
}

func TestCallback(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abcdef", nil)
	w := httptest.NewRecorder()
	callBackFunc()(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestServer(t *testing.T) {
	ch := make(chan os.Signal)
	go func() {
		Start(8080)
		<-ch
	}()
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
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
