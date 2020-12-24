package web

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, out interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(out); err != nil {
		return err
	}
	return nil
}
