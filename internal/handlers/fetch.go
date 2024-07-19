package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var query = "SELECT * FROM `hack-team-tenali-rama.unwind_events.life-events` ORDER BY Time DESC LIMIT 5 -- %s %s"

func Fetch(w http.ResponseWriter, r *http.Request, client *bigquery.Client) ([]map[string]interface{}, error) {
	ctx := context.Background()
	start, end := getWeekStartAndEnd()
	q := fmt.Sprintf(query, start.Format(time.RFC3339), end.Format(time.RFC3339))
	res := client.Query(q)
	it, err := res.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	var entries []map[string]interface{}
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate through results: %w", err)
		}

		// Convert BigQuery values to Go values
		entry := make(map[string]interface{})
		for column, value := range row {
			entry[column] = value
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func getWeekStartAndEnd() (time.Time, time.Time) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	return startOfWeek, endOfWeek
}
