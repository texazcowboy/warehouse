package web

import (
	"net/http"
	"strings"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"
)

type TokenValidator func(string) error

func AuthenticationMiddleware(next http.HandlerFunc, l *logger.LogEntryWithStackTrace, validatorFunc TokenValidator) http.HandlerFunc { // nolint
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		arr := strings.Split(authHeader, " ")
		if len(arr) != 2 { // nolint
			if err := RespondError(w, http.StatusBadRequest, "unexpected tokens in authorization header"); err != nil {
				l.WithError(err).Error("Unable to write response")
			}
			return
		}
		if arr[0] != "Bearer" {
			if err := RespondError(w, http.StatusBadRequest, "unsupported authorization type"); err != nil {
				l.WithError(err).Error("Unable to write response")
			}
			return
		}
		err := validatorFunc(arr[1])
		if err != nil {
			if err := RespondError(w, http.StatusUnauthorized, "unauthorized"); err != nil {
				l.WithError(err).Error("Unable to write response")
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
