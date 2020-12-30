package web

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type TokenValidator func(string) error

func AuthenticationMiddleware(next http.HandlerFunc, l *logrus.Entry, validatorFunc TokenValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		arr := strings.Split(authHeader, " ")
		if len(arr) != 2 {
			if err := RespondError(w, http.StatusBadRequest, "unexpected tokens in authorization header"); err != nil {
				l.Error(err)
			}
			return
		}
		err := validatorFunc(arr[1])
		if err != nil {
			if err := RespondError(w, http.StatusUnauthorized, "unauthorized"); err != nil {
				l.Error(err)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
