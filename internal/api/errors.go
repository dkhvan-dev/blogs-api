package api

import (
	"blogs-api/internal/core/errors"
	"encoding/json"
	"net/http"
)

func handleError(w http.ResponseWriter, err errors.Error) {
	errJSON := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	data, marshalErr := json.Marshal(errJSON)
	if marshalErr != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	_, writeErr := w.Write(data)
	if writeErr != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	return
}
