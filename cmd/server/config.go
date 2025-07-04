package main

import (
    "log/slog"
    "os"
    "strconv"
    "time"
)

type config struct {
    port        string
    logLevel    slog.Level
    shutdownDur time.Duration
}

func newConfig() (*config, error) {
    lvlStr := os.Getenv("LOG_LEVEL")
    lvl := slog.LevelInfo
    if lvlStr != "" {
        _ = lvl.UnmarshalText([]byte(lvlStr))
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    shutdown := 15 * time.Second
    if s := os.Getenv("SHUTDOWN_TIMEOUT_SEC"); s != "" {
        if v, err := strconv.Atoi(s); err == nil {
            shutdown = time.Duration(v) * time.Second
        }
    }

    return &config{
        port:        port,
        logLevel:    lvl,
        shutdownDur: shutdown,
    }, nil
}
