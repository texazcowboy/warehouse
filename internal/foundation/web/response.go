package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorResponse struct {
	Message string
}

func Respond(w http.ResponseWriter, status int, payload interface{}) error {
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return nil
	}

	response, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "Respond -> json.Marshal(%v)", payload)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "Respond -> w.Write(%v)", response)
	}

	return nil
}

func RespondError(w http.ResponseWriter, status int, message string) error {
	return Respond(w, status, &ErrorResponse{message})
}
