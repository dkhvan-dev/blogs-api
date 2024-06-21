package api

import (
	"encoding/json"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	errJSON := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	data, err := json.Marshal(errJSON)
	if err != nil {
		data = []byte(`{"error": "internal server error"}`)
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
