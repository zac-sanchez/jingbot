package web

import (
	"context"

	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func New(logger *zap.Logger, api *slack.Client, schedule string, randomness int) *mux.Router {

	backend := &API{logger, api, schedule, randomness}
	frontend := Frontend(
		&HomepageData{
			"JingBot",
			[]string{"/api/v1/hello", "/api/v1/schedule", "/api/v1/channels"},
		},
		"static/home.html.tmpl",
		)

	r := mux.NewRouter()
	r.Handle("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`OK`))
	}))
	r.HandleFunc("/api/v1/hello", backend.Hello)
	r.HandleFunc("/api/v1/schedule", backend.Schedule)
	r.HandleFunc("/api/v1/channels", backend.Channels)
	r.HandleFunc("/", frontend)
	r.HandleFunc("/home", frontend)
	return r
}

func Run(r *mux.Router, logger *zap.Logger, port string, ctx context.Context) {
	s := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down web server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()

	logger.Info("Listening", zap.String("port", ":"+port))
	if err := s.ListenAndServe(); err != nil {
		logger.Error("Web server terminated: %v", zap.Error(err))
	}
}
