package main

import (
    "io"
    "log/slog"
)

func newLogger(w io.Writer, lvl slog.Level) *slog.Logger {
    return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: lvl}))
}
