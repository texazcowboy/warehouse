package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Respond(w http.ResponseWriter, status int, payload interface{}) error {
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return nil
	}

	response, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "web: Respond -> json.Marshal(%v)", payload)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "web: Respond -> w.Write(%v)", response)
	}

	return nil
}

func RespondError(w http.ResponseWriter, status int, message string) error {
	return Respond(w, status, &ErrorResponse{message})
}
