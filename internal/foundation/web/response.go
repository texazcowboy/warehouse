package web

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string
}

func Respond(w http.ResponseWriter, status int, payload interface{}) error {
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func ResponseError(w http.ResponseWriter, status int, message string) error {
	return Respond(w, status, &ErrorResponse{message})
}
