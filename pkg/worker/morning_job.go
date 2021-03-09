package worker

import (
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

// Returns a function that posts a message as bot user
func Morning(logger *zap.Logger, api *slack.Client, minutes int) func() {
	return func() {
		// TODO: Avoid running this every iteration by updating the store when bot is added to a new channel
		channels, _, err := api.GetConversationsForUser(&slack.GetConversationsForUserParameters{})
		if err != nil {
			zap.Error(err)
			panic(err)
		}
		r := rand.Intn(minutes)
		time.Sleep(time.Duration(r) * time.Minute)
		for _, c := range channels {
			// Randomly wait some time between 0 - 10 minutes to post the update
			channelID, timestamp, err := api.PostMessage(c.ID, slack.MsgOptionText(":remote-sleepy-morning:", false))
			if err != nil {
				logger.Error("Error posting message",
					zap.String("channel", channelID),
					zap.String("timestamp", timestamp),
					zap.Error(err),
				)
				continue
			}
			logger.Info("Posted message", zap.String("channel", channelID), zap.String("timestamp", timestamp))
		}
	}
}
