package security

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type Validator func(string) error

func AuthenticationMiddleware(next http.HandlerFunc, l *logrus.Entry, validatorFunc Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		arr := strings.Split(authHeader, " ")
		if len(arr) != 2 {
			if err := web.RespondError(w, http.StatusBadRequest, "unauthorized"); err != nil {
				l.Error(err)
			}
			return
		}
		err := validatorFunc(arr[1])
		if err != nil {
			if err := web.RespondError(w, http.StatusUnauthorized, "unauthorized"); err != nil {
				l.Error(err)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
