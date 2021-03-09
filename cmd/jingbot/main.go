package main

import (
	"context"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
	"go.uber.org/zap"

	"bitbucket.org/zac_sanchez/jingbot/pkg/web"
	"bitbucket.org/zac_sanchez/jingbot/pkg/worker"
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
	_, err := c.AddFunc(cfg.Time, worker.Morning(logger, api, cfg.Minutes))
	if err != nil {
		zap.Error(err)
		panic(err)
	}
	c.Start()

	r := web.New(logger, api, cfg.Time, cfg.Minutes)
	go web.Run(r, logger, cfg.Port, ctx)

	<-ctx.Done()
	c.Stop()
	logger.Info("Received termination signal. Waiting before exit.")
	time.Sleep(time.Second)
	logger.Info("Exiting")

}
