package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func ExtractIntFromRequest(r *http.Request, varName string) (int, error) {
	vars := mux.Vars(r)
	pathVar := vars[varName]
	id, err := strconv.Atoi(pathVar)
	if err != nil {
		return 0, errors.Wrapf(err, "web: ExtractIntFromRequest -> strconv.Atoi(%v)", pathVar)
	}
	return id, nil
}
