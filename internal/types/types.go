package types

import (
	"time"
)

type Event struct {
	SuppliedBy       string    `json:"suppliedBy" bigquery:"suppliedBy"`
	SuppliedFor      string    `json:"suppliedFor" bigquery:"suppliedFor"`
	Time             time.Time `json:"time" bigquery:"time"`
	EventName        string    `json:"eventName" bigquery:"eventName"`
	PeopleInvolved   []string  `json:"peopleInvolved" bigquery:"peopleInvolved"`
	EventType        string    `json:"eventType" bigquery:"eventType"`
	Activities       []string  `json:"activities" bigquery:"activities"`
	Vibe             string    `json:"vibe" bigquery:"vibe"`
	ThingsToRemember []string  `json:"thingsToRemember"  bigquery:"thingsToRemember"`
}

type UserInfo struct {
	Username string
}

type Info struct {
	SuppliedBy  string    `json:"suppliedBy"`
	SuppliedFor string    `json:"suppliedFor"`
	Time        time.Time `json:"time"`
	Text        string    `json:"text"`
}

type TimelineEntry struct {
	Title   string
	Content string
	Date    time.Time
}
