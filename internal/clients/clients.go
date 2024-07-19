package clients

import (
	"context"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
)

func GetVertexClient() (*genai.Client, error) {
	ctx := context.Background()
	return genai.NewClient(ctx, "hack-team-tenali-rama", "us-central1")
}

func GetBigQueryClient() (*bigquery.Client, error) {
	ctx := context.Background()
	return bigquery.NewClient(ctx, "hack-team-tenali-rama")
}

func GetStorageClient() (*storage.Client, error) {
	ctx := context.Background()
	return storage.NewClient(ctx)
}
