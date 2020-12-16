package main

import (
	"context"
	"github.com/sparfenov/httpmux/internal/app"
	"github.com/sparfenov/httpmux/internal/server/handlers"
	"github.com/sparfenov/httpmux/internal/server/middlewares"
	"github.com/sparfenov/httpmux/pkg/httpserver"
	"github.com/sparfenov/httpmux/pkg/logger"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	config, err := app.NewConfig()
	if err != nil {
		panic("failed to create app, configuration error: " + err.Error())
	}

	l := logger.NewLogger(config.IsDebug)

	go func() {
		sig := <-sigChan
		l.Infof("received %s os signal", sig)

		cancel()
	}()

	serveHTTP(ctx, l, config)
}

func serveHTTP(ctx context.Context, l logger.Interface, config *app.Config) {
	srv := httpserver.Server{
		Server: &http.Server{
			Addr: config.Addr,
			Handler: handlers.URLMuxHandler{
				Logger:                 l,
				MaxURLCountToProcess:   config.MaxURLCountToProcess,
				ExternalRequestLimit:   config.ExternalRequestLimit,
				ExternalRequestTimeout: config.ExternalRequestTimeout,
			},
		},
		LimitReq: config.ListenLimit,
	}

	srv.Use(middlewares.NewLoggingMiddleware(l).LoggingMiddleware)

	go func() {
		l.Infof("Starting server on: %s", config.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("listen err: %s", err)
		}
	}()

	<-ctx.Done()

	l.Infof("Server is shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ServerShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Errorf("Shutdown err: %s", err)
	}

	l.Infof("Server successfully shut down")
}
