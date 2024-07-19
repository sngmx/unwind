package router

import (
	"net/http"
	"unwind/internal/handlers"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
)

func NewRouter(vertexClient *genai.Client, bigQueryClient *bigquery.Client, storageClient *storage.Client) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/timeline", func(w http.ResponseWriter, r *http.Request) {
		handlers.Timeline(w, r, bigQueryClient)
	})
	mux.HandleFunc("/upload", handlers.Upload)
	mux.HandleFunc("/submit-username", handlers.SubmitUsername)
	mux.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		handlers.Extract(w, r, vertexClient, bigQueryClient, storageClient)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}
