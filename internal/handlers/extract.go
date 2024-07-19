package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"unwind/internal/cookies"
	"unwind/internal/models" // Import the models package
	"unwind/internal/types"
	"unwind/internal/utils"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
)

var (
	datasetID  = "unwind_events"
	tableName  = "life-events"
	bucketName = "unwind-storage"
)

func Extract(w http.ResponseWriter, r *http.Request, client *genai.Client, bigQueryClient *bigquery.Client, storageClient *storage.Client) {
	ctx := context.Background()
	generativeModel := models.BuildVertexModel(client)
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		user, _ := utils.GetUserInfo(w, r, cookies.Store)
		textData := r.FormValue("textData")

		file, fileHeader, err := r.FormFile("fileData")
		if err != nil {
			log.Println("File not provided")
		} else {
			uri, err := syncFile(ctx, file, fileHeader, storageClient)
			if err != nil {
				log.Fatal("Cannot upload file.", err)
			} else {
				info := &types.Info{
					SuppliedBy:  user.Username,
					SuppliedFor: "Sangam",
					Time:        time.Now(),
				}
				jsonText, _ := json.Marshal(info)
				resp, err := generativeModel.GenerateContent(ctx, genai.FileData{
					FileURI:  uri,
					MIMEType: fileHeader.Header.Get("Content-Type"),
				}, genai.Text(jsonText))
				sendToBQ(ctx, err, w, resp, bigQueryClient)
			}
			defer file.Close()
		}
		if textData != "" {
			info := &types.Info{
				SuppliedBy:  user.Username,
				SuppliedFor: "Sangam",
				Time:        time.Now(),
				Text:        textData,
			}
			jsonText, _ := json.Marshal(info)
			resp, err := generativeModel.GenerateContent(ctx, genai.Text(jsonText))
			sendToBQ(ctx, err, w, resp, bigQueryClient)
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	fmt.Fprintf(w, "Success")
}

func sendToBQ(ctx context.Context, err error, w http.ResponseWriter, resp *genai.GenerateContentResponse, bigQueryClient *bigquery.Client) {
	if err != nil {
		log.Fatal("error generating content: %w", err)
		w.Write([]byte("Service unavailable"))
	}
	events := parse(resp)
	if events != nil {
		if err = save(ctx, events, bigQueryClient); err != nil {
			log.Fatal("Error saving data to big query", err)
		}
	}
}

func parse(resp *genai.GenerateContentResponse) []types.Event {
	var events []types.Event
	if len(resp.Candidates) > 0 {
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				fmt.Println(part.(genai.Text))
				input := part.(genai.Text)
				cleanedInput := strings.ReplaceAll(string(input), "```json", "")
				cleanedInput = strings.ReplaceAll(cleanedInput, "```", "")
				cleanedInput = strings.TrimSpace(cleanedInput)
				jsonerror := json.Unmarshal([]byte(cleanedInput), &events)
				if jsonerror != nil {
					fmt.Println("Error unmarshaling JSON:", jsonerror)
					return nil
				} else {
					return events
				}
			}
		}
	} else {
		fmt.Println("No candidates found in the response.")
	}
	return nil
}

func save(ctx context.Context, events []types.Event, client *bigquery.Client) error {
	table := client.Dataset(datasetID).Table(tableName)
	u := table.Uploader()

	items := make([]*bigquery.StructSaver, len(events))
	for i, p := range events {
		items[i] = &bigquery.StructSaver{
			Struct: p,
		}
	}
	if err := u.Put(ctx, items); err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

func syncFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, client *storage.Client) (string, error) {
	bucket := client.Bucket(bucketName)

	fileName := header.Filename
	wc := bucket.Object(fileName).NewWriter(ctx)
	_, err := io.Copy(wc, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	return fmt.Sprintf("gs://%s/%s", bucketName, fileName), nil
}
