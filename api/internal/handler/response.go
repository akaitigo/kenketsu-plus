package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// M-7: classify errors — ValidationError→400, everything else→500
func writeRepoError(w http.ResponseWriter, err error) {
	var validationErr *model.ValidationError
	if errors.As(err, &validationErr) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeError(w, http.StatusInternalServerError, "internal server error")
}
