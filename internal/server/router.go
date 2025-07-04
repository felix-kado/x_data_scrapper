package server

import (
    "net/http"

    "github.com/go-chi/chi"
    "github.com/felix-kado/x_data_scrapper/internal/handler"
)

// Собирает основной роутер chi
func NewRouter(user *handler.UserHandler, expand *handler.ExpandHandler, metrics *handler.MetricsHandler) http.Handler {
    r := chi.NewRouter()

    r.Mount("/users", user.Routes())
    r.Mount("/metrics", metrics)
    r.Post("/expand", expand.ServeHTTP)

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    return r
}
