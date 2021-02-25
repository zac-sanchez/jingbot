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

	channels, err := getBotChannels(cfg.BotID, api)
	if err != nil {
		zap.Error(err)
		panic(err)
	}

	c := cron.New()
	_, err = c.AddFunc(cfg.Time, jingFunc(logger, channels, api))
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

// Returns a function that posts a message as Jing
func jingFunc(logger *zap.Logger, channels []slack.Channel, api *slack.Client) func() {
	return func() {
		for _, c := range channels {
			channelID, timestamp, err := api.PostMessage(c.ID, slack.MsgOptionText(":remote-sleepy-morning:", false))
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Minute)
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


var blocklistedChannels = []string{
	"general",
	"announcements",
}

func filterChannel(ss []slack.Channel, test func(string) bool) (ret []slack.Channel) {
	for _, s := range ss {
		if test(s.Name) {
			ret = append(ret, s)
		}
	}
	return
}

func notBlocked(s string) bool {
	for _, blocked := range blocklistedChannels {
		if s == blocked {
			return false
		}
	}
	return true
}

// post only to channels other than general and announcement where JingBot is present.
func getBotChannels(botID string, api *slack.Client) ([]slack.Channel, error) {
	channels, _, err := api.GetConversationsForUser(&slack.GetConversationsForUserParameters{
		UserID: botID,
	})
	if err != nil {
		return []slack.Channel{}, err
	}
	channels = filterChannel(channels, notBlocked)
	return channels, nil
}
