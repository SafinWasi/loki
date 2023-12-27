package web

import (
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

func ParseTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/*.html"))
}
func Start(port int) {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("/", homeHandler())
	mux.Handle("/registration", registrationHandler())
	mux.Handle("/delete/", deleteHandler())
	mux.Handle("/client/", clientHandler())
	mux.Handle("/callback", callBackFunc())
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
		err = tp.ExecuteTemplate(w, "client", string(val))
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func deleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNoContent)
	}
}

func registrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			host := r.FormValue("host")
			payload := r.FormValue("payload")
			if payload == "" {
				var test openid.RegistrationPayload
				test.Ssa = r.FormValue("ssa")
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

	}
}

func callBackFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Query())
	}
}
