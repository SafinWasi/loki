package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/safinwasi/loki/openid"
	"github.com/safinwasi/loki/secrets"
)

var tp = ParseTemplates()

func ParseTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/*.html"))
}
func Start(port int) {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	keys, _ := secrets.GetKeys()
	mux.Handle("/", homeHandler(keys))
	mux.Handle("/registration", registrationHandler(nil))
	fmt.Printf("Starting Loki on http://127.0.0.1:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func homeHandler(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tp.ExecuteTemplate(w, "home", data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}

func registrationHandler(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			host := r.FormValue("host")
			ssa := r.FormValue("SSA")
			scope := r.FormValue("scope")
			_ = scope
			newClient, err := openid.Register(host, ssa, "")
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
			secrets.Set(host, clientBytes)
		} else {
			err := tp.ExecuteTemplate(w, "registration", data)
			if err != nil {
				log.Println(err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
	}
}
