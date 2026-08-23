package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Alex-bondaruk/rariable/internal/config"
	"github.com/Alex-bondaruk/rariable/internal/rarible"
	"github.com/Alex-bondaruk/rariable/internal/service"
)

func main() {
	cfgLog := zap.NewProductionConfig()
	cfgLog.DisableStacktrace = true
	log, err := cfgLog.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config", zap.Error(err))
	}

	client, err := rarible.New(cfg.BaseURL, cfg.APIKey)
	if err != nil {
		log.Fatal("rarible client", zap.Error(err))
	}

	router := service.NewRouter(service.NewHandler(client, log))
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("shutdown", zap.Error(err))
	}
}
