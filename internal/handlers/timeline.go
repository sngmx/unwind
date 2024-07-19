package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"unwind/internal/cookies"
	"unwind/internal/utils"
)

func Timeline(w http.ResponseWriter, r *http.Request) {

	user, _ := utils.GetUserInfo(w, r, cookies.Store)
	path := filepath.Join("internal", "templates", "timeline.html")
	tmpl := template.Must(template.ParseFiles(path))
	tmpl.Execute(w, user.Username)
}
