package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"


    "github.com/felix-kado/x_data_scrapper/internal/handler"
    "github.com/felix-kado/x_data_scrapper/internal/server"
    "github.com/felix-kado/x_data_scrapper/internal/service"
    "github.com/felix-kado/x_data_scrapper/internal/twitter"
)

func main() {
    cfg, err := newConfig()
    if err != nil {
        panic(err)
    }

    logger := newLogger(os.Stdout, cfg.logLevel)

    bearer := os.Getenv("TWITTER_BEARER_TOKEN")
    if bearer == "" {
        logger.Error("TWITTER_BEARER_TOKEN env var required")
        os.Exit(1)
    }

    twClient := twitter.NewClient(nil, bearer)

    userSvc := service.NewUserService(twClient)
    expandSvc := service.NewExpandService(twClient, userSvc)

    // handlers
    userHandler := handler.NewUserHandler(userSvc)
    expandHandler := handler.NewExpandHandler(expandSvc)
    metricsHandler := handler.NewMetricsHandler(userSvc)

    mux := server.NewRouter(userHandler, expandHandler, metricsHandler)

    srv := &http.Server{
        Addr:    ":" + cfg.port,
        Handler: mux,
    }

    go func() {
        logger.Info("http server listening", slog.String("addr", srv.Addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server", slog.Any("error", err))
            os.Exit(1)
        }
    }()

    // graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), cfg.shutdownDur)
    defer cancel()

    _ = srv.Shutdown(ctx)
}
