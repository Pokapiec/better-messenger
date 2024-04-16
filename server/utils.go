package main

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	data interface{}
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
