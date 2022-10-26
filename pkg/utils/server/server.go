package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

func RunHTTP(ctx context.Context, s *http.Server, quit chan<- struct{}) *http.Server {
	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.WithError(log.FromContext(ctx), err).Error("http server start error")
			}
			quit <- struct{}{}
		}
	}()

	return s
}

func ShutdownHTTP(ctx context.Context, srv *http.Server) {
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(log.FromContext(ctx), err).Error("unable to shutdown the http server")
	} else {
		log.FromContext(ctx).Info("http server shutdown success")
	}
}

func GracefulShutdown(ctx context.Context, quit <-chan struct{}, cleanup func(ctx context.Context)) {
	const shutdownWaitingTime = 10 * time.Second

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.FromContext(ctx).Info("received shutdown signal")
	case <-quit:
		log.FromContext(ctx).Info("received quit notification")
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownWaitingTime)
	defer cancel()

	cleanup(ctx)

	log.FromContext(ctx).Info("server shutdown success")
}
