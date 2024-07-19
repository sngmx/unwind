package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"unwind/internal/cookies"
)

func Home(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("internal", "templates", "home.html")
	tmpl := template.Must(template.ParseFiles(path))
	tmpl.Execute(w, nil)
}

func SubmitUsername(w http.ResponseWriter, r *http.Request) {
	session, err := cookies.Store.Get(r, "unwind-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/timeline", http.StatusFound)
}
