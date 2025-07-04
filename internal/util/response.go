package util

import (
    "encoding/json"
    "net/http"
)

type apiError struct {
    Error string `json:"error"`
}

func JSON(w http.ResponseWriter, err error, data any) {
    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(apiError{Error: err.Error()})
        return
    }
    _ = json.NewEncoder(w).Encode(data)
}
