package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"my-calendar/internal/calendar"
	"my-calendar/internal/clickhouse"
	"my-calendar/internal/config"
	"my-calendar/internal/handler"
	_ "my-calendar/internal/handler/docs"
	"my-calendar/internal/logger"
	"my-calendar/internal/metrics"
	"my-calendar/internal/middleware"
)

// @title Calendar API
// @version 1.0
// @description REST API для календаря
// @host localhost:8080
// @BasePath /
func main() {

	// first initialization of logger
	logger.Logger()
	defer func() {
		_ = logger.L().Sync()
	}()

	if err := clickhouse.ConnectWithRetry(); err != nil {
		logger.L().Warn(
			"clickhouse unavailable",
			zap.Error(err),
		)
	}

	// load config
	cfg := config.LoadCfg()

	// var for event storage
	var calendar = calendar.NewCalendar()

	// var for handler
	h := handler.NewHandler(calendar)

	// create router
	r := mux.NewRouter()
	r.Use(middleware.RequestLogger)
	r.HandleFunc("/events", h.GetEvents).Methods("get")
	r.HandleFunc("/events", h.CreateEvent).Methods("post")
	r.HandleFunc("/events/{id}", h.UpdateEvent).Methods("put")
	r.HandleFunc("/events/{id}", h.DeleteEvent).Methods("delete")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Handle("/metrics", promhttp.Handler())

	// HTTP Server setup
	srv := &http.Server{
		Addr:              cfg.HttpPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.L().Info("starting server",
		zap.String("host", "localhost"),
		zap.String("port", cfg.HttpPort),
	)

	metrics.Init()

	// create waitgroup
	wg := &sync.WaitGroup{}

	// run http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.L().Info("http server starting", zap.String("addr", "localhost"+srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal("http server failed", zap.Error(err))
		}
	}()

	// signal channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// wait for signal
	sig := <-sigCh
	logger.L().Info("shutdown signal received", zap.String("signal", sig.String()))

	// graceful HTTP shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.L().Error("http server shutdown error", zap.Error(err))
	} else {
		logger.L().Info("http server stopped gracefully")
	}

	// wait for goroutines to finish
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		logger.L().Info("all goroutines finished")
	case <-time.After(20 * time.Second):
		logger.L().Warn("timeout waiting for goroutines, forcing exit")
	}

	logger.L().Info("service stopped")
}
