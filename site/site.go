package site

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type PageData struct {
	URLs []string
}

func StartSiteProcess(config *Config) {
	r := chi.NewRouter()

	// Parse all templates
	tmpl, err := template.ParseGlob("site/templates/*.tmpl.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			URLs: config.URLs,
		}
		if err := tmpl.ExecuteTemplate(w, "home.tmpl.html", data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		}
	})

	log.Println("Starting server on localhost:3451")
	http.ListenAndServe(":3451", r)

}
