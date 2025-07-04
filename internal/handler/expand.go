package handler

import (
    "encoding/json"
    "net/http"

    "github.com/felix-kado/x_data_scrapper/internal/service"
    "github.com/felix-kado/x_data_scrapper/internal/util"
)

// HTTP-обработчик /expand
 type ExpandHandler struct {
    svc *service.ExpandService
}

func NewExpandHandler(svc *service.ExpandService) *ExpandHandler {
    return &ExpandHandler{svc: svc}
}

func (h *ExpandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var p service.ExpandParams
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        util.JSON(w, err, nil)
        return
    }

    graph, err := h.svc.Expand(r.Context(), p)
    util.JSON(w, err, graph)
}
