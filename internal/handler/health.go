package handler

import (
	"encoding/json"
	"net/http"
)

type HeathResponse struct {
	Status string `json:"status"`
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := HeathResponse{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
