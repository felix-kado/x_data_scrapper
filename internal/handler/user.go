package handler

import (
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
    "github.com/felix-kado/x_data_scrapper/internal/service"
    "github.com/felix-kado/x_data_scrapper/internal/util"
)

// HTTP-обработчик пользователей
 type UserHandler struct {
    svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{svc: svc}
}

func (h *UserHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/{username}", h.getProfile)
    r.Get("/{username}/tweets", h.getTweets)
    r.Get("/{username}/metrics", h.getMetrics)
    return r
}

func (h *UserHandler) getProfile(w http.ResponseWriter, r *http.Request) {
    username := chi.URLParam(r, "username")
    ctx := r.Context()

    user, err := h.svc.GetProfile(ctx, username)
    util.JSON(w, err, user)
}

func (h *UserHandler) getTweets(w http.ResponseWriter, r *http.Request) {
    username := chi.URLParam(r, "username")
    limitStr := r.URL.Query().Get("limit")
    limit := 0
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }
    tweets, err := h.svc.GetTweets(r.Context(), username, limit)
    util.JSON(w, err, tweets)
}

func (h *UserHandler) getMetrics(w http.ResponseWriter, r *http.Request) {
    username := chi.URLParam(r, "username")
    metrics, err := h.svc.ComputeMetrics(r.Context(), username, 200) // sample size
    util.JSON(w, err, metrics)
}
