package main

import (
	"log"
	"net/http"
	"unwind/internal/clients"
	"unwind/internal/router"
)

func main() {

	var err error
	vertexClient, err := clients.GetVertexClient()
	if err != nil {
		log.Fatalf("Failed to initialize Vertex AI: %v", err)
	}

	bigQueryClient, err := clients.GetBigQueryClient()
	if err != nil {
		log.Fatalf("Failed to initialize BigQueryClient: %v", err)
	}

	storageClient, err := clients.GetStorageClient()
	if err != nil {
		log.Fatalf("Failed to initialize StorageClient: %v", err)
	}
	r := router.NewRouter(vertexClient, bigQueryClient, storageClient)

	log.Println("Server is running at ::8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
