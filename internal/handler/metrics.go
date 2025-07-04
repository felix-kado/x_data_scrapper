package handler

import (
    "net/http"

    "github.com/go-chi/chi"
    "github.com/felix-kado/x_data_scrapper/internal/service"
    "github.com/felix-kado/x_data_scrapper/internal/util"
)

// HTTP-обработчик /metrics
 type MetricsHandler struct {
    svc *service.UserService
}

func NewMetricsHandler(svc *service.UserService) *MetricsHandler {
    return &MetricsHandler{svc: svc}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    username := chi.URLParam(r, "username")
    m, err := h.svc.ComputeMetrics(r.Context(), username, 200)
    util.JSON(w, err, m)
}
