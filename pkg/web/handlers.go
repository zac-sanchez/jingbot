package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type API struct {
	logger     *zap.Logger
	slack      *slack.Client
	schedule   string
	randomness int
}

type HelloPayload struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type SchedulePayload struct {
	Randomness     string `json:"random_interval"`
	CronExpression string `json:"cron_expression"`
}

func (api *API) Channels(w http.ResponseWriter, r *http.Request) {
	channels, _, err := api.slack.GetConversationsForUser(&slack.GetConversationsForUserParameters{})
	if err != nil {
		http.Error(w, "Error getting conversations from Slack API", http.StatusInternalServerError)
	}

	var names []struct {
		Channel string `json:"channel"`
	}
	for _, c := range channels {
		names = append(names, struct {
			Channel string `json:"channel"`
		}{c.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(names); err != nil {
		http.Error(w, "Unable to construct message", http.StatusInternalServerError)
	}
}

func (api *API) Hello(w http.ResponseWriter, req *http.Request) {
	payload := HelloPayload{
		"in_channel",
		":remote-sleepy-morning:",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Unable to construct message", http.StatusInternalServerError)
	}

}

func (api *API) Schedule(w http.ResponseWriter, req *http.Request) {

	payload := SchedulePayload{
		strconv.Itoa(api.randomness),
		api.schedule,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Unable to construct message", http.StatusInternalServerError)
	}

}
