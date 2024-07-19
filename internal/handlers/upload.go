package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("internal", "templates", "upload.html")
	tmpl := template.Must(template.ParseFiles(path))
	tmpl.Execute(w, nil)
}
