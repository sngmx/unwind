package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"unwind/internal/cookies"
	"unwind/internal/utils"

	"cloud.google.com/go/bigquery"
)

func Timeline(w http.ResponseWriter, r *http.Request, client *bigquery.Client) {

	user, _ := utils.GetUserInfo(w, r, cookies.Store)

	entries, _ := Fetch(w, r, client) // Replace 'handlers.Fetch' with 'Fetch' and discard the error
	timeline_tabs := filepath.Join("internal", "templates", "timeline_tabs.html")
	timeline := filepath.Join("internal", "templates", "timeline.html")

	funcMap := template.FuncMap{
		"formatValue": formatValue,
	}

	// Parse the templates
	tmpl, err := template.New("timeline_tabs.html").Funcs(funcMap).ParseFiles(timeline, timeline_tabs)
	if err != nil {
		http.Error(w, "Failed to parse templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the layout template with the context
	formattedEntries := make([]map[string]string, len(entries))
	for i, entry := range entries {
		formattedEntry := make(map[string]string)
		for k, v := range entry {
			formattedValue := formatValue(v)
			if formattedValue != "" || (k == "Time" || k == "EventName" || k == "Vibe") {
				formattedEntry[k] = formattedValue
			}
		}
		formattedEntries[i] = formattedEntry
	}
	err = tmpl.ExecuteTemplate(w, "timeline_tabs.html", struct {
		Entries  []map[string]string
		Username string
	}{
		Entries:  formattedEntries,
		Username: user.Username,
	})
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case []bigquery.Value: // Handle slices
		if len(v) == 0 {
			return "" // Return empty string for empty lists
		}
		var strSlice string
		for _, elem := range v {
			strSlice += formatValue(elem) + ","
		}
		return strSlice // Join with commas
	default:
		return fmt.Sprint(value)
	}
}
