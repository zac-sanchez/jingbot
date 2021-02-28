package main

import (
	"bitbucket.org/zac_sanchez/jingbot/pkg/web"
	"context"
	"math/rand"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func main() {

	ctx := WithTerminationSignals(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	cfg := getConfig()
	api := slack.New(cfg.APIKey)

	logger.Info("Started JingBot",
		zap.String("Environment", cfg.Environment),
		zap.String("cron string", cfg.Time),
		zap.String("port", cfg.Port),
		zap.Int("randomness interval in minutes", cfg.Minutes),
	)

	c := cron.New()
	_, err := c.AddFunc(cfg.Time, worker(logger, api, cfg.Minutes))
	if err != nil {
		zap.Error(err)
		panic(err)
	}
	c.Start()

	r := web.New(logger)
	go web.Run(r, logger, cfg.Port, ctx)

	<-ctx.Done()
	c.Stop()
	logger.Info("Received termination signal. Waiting before exit.")
	time.Sleep(time.Second)
	logger.Info("Exiting")

}

// Returns a function that posts a message as bot user
func worker(logger *zap.Logger, api *slack.Client, minutes int) func() {
	return func() {
		// TODO: Avoid running this every iteration by updating the store when bot is added to a new channel
		channels, _, err := api.GetConversationsForUser(&slack.GetConversationsForUserParameters{})
		if err != nil {
			zap.Error(err)
			panic(err)
		}
		for _, c := range channels {
			// Randomly wait some time between 0 - 10 minutes to post the update
			r := rand.Intn(minutes)
			time.Sleep(time.Duration(r) * time.Minute)
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