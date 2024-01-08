package web

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/99designs/keyring"
	"github.com/safinwasi/loki/openid"
	"github.com/safinwasi/loki/secrets"
)

var tp = ParseTemplates()
var currentOP openid.Configuration

//go:embed html/*.html
//go:embed static/*
var content embed.FS

func ParseTemplates() *template.Template {
	return template.Must(template.ParseFS(content, "html/*.html"))
}
func Start(port int) {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/", http.FileServer(http.FS(content))))
	mux.Handle("/", homeHandler())
	mux.Handle("/registration", registrationHandler())
	mux.Handle("/delete/", deleteHandler())
	mux.Handle("/client/", clientHandler())
	mux.Handle("/callback", callBackFunc())
	mux.Handle("/code/", codeFlow())
	mux.Handle("/add", addFunc())
	fmt.Printf("Starting Loki on http://127.0.0.1:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := secrets.GetKeys()
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		err = tp.ExecuteTemplate(w, "home", keys)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}

func clientHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/client/")
		val, err := secrets.Get(id)
		if err != nil {
			log.Println(err)
			if errors.Is(err, keyring.ErrKeyNotFound) {
				http.Error(w, fmt.Sprintf("%s not found", id), http.StatusNotFound)
				return
			} else {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
		_, err = w.Write(val)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func deleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Endpoint expects DELETE", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/delete/")
		_, err := secrets.Get(id)
		if err != nil {
			log.Println(err)
			if errors.Is(err, keyring.ErrKeyNotFound) {
				http.Error(w, fmt.Sprintf("%s not found", id), http.StatusNotFound)
				return
			} else {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
		err = secrets.RemoveKey(id)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully removed %s\n", id)
		w.Write([]byte(""))
	}
}

func registrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			host := r.FormValue("host")
			payload := r.FormValue("payload")
			log.Println("Host and payload:", host, payload)
			if payload == "" {
				var test openid.RegistrationPayload
				test.Ssa = r.FormValue("ssa")
				test.Grant_types = []string{"code", "client_credentials"}
				test.Client_name = r.FormValue("client_name")
				b, err := json.Marshal(test)
				if err != nil {
					log.Println(err)
					return
				}
				payload = string(b)
			}
			newClient, err := openid.Register(host, "")
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			clientBytes, err := json.Marshal(newClient)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			hostName, err := url.Parse(host)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			secrets.Set(hostName.Host, clientBytes)
			w.Write([]byte("<p>Successfully registered</p>"))
		} else {
			err := tp.ExecuteTemplate(w, "registration", nil)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
	}
}

func codeFlow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/code/")
		if len(id) == 0 {
			return
		}
		val, err := secrets.Get(id)
		if err != nil {
			log.Println(err)
			if errors.Is(err, keyring.ErrKeyNotFound) {
				http.Error(w, fmt.Sprintf("%s not found", id), http.StatusNotFound)
				return
			} else {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
		var config openid.Configuration
		json.Unmarshal(val, &config)
		uri := CreateCodeUrl(config)
		currentOP = config
		uri = fmt.Sprintf("%s?%s", config.OpenID.Authorization_endpoint, uri)
		http.Redirect(w, r, uri, http.StatusFound)
	}
}

func callBackFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		code := params.Get("code")
		log.Println(code)
		token, err := SendTokenRequest(code, currentOP.Client_id, currentOP.Client_secret, currentOP.OpenID.Token_endpoint, "authorization_code")
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		tp.ExecuteTemplate(w, "callback", token)
	}
}

func addFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Endpoint expects POST", http.StatusMethodNotAllowed)
			return
		}
		var newClient openid.Configuration
		r.ParseForm()
		newClient.Client_id = r.FormValue("client_id")
		newClient.Client_secret = r.FormValue("client_secret")
		host := r.FormValue("hostname")
		hostName, err := url.Parse(host)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		oidc, err := openid.Fetch_openid(host)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		newClient.OpenID = *oidc
		clientBytes, err := json.MarshalIndent(newClient, "", "\t")
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		secrets.Set(hostName.Host, clientBytes)
		w.Write([]byte("<p>Successfully added</p>"))
	}
}
