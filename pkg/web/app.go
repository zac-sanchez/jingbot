package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"

	_ "bitbucket.org/zac_sanchez/jingbot/pkg/web/statik"
)

func New(logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()
	sfs, err := fs.New()
	if err != nil {
		panic("Failed to create frontend")
	}
	r.Handle("/", http.FileServer(sfs))
	r.Handle("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`OK`))
	}))
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

	logger.Info("Listening", zap.String("port", ":" + port))
	if err := s.ListenAndServe(); err != nil {
		logger.Error("Web server terminated: %v", zap.Error(err))
	}
}