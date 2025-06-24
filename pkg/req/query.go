package req

import (
	"errors"
	"net/http"
)

func GetQueryParam(r *http.Request, name string) (string, error) {
	val := r.URL.Query().Get(name)
	if val == "" {
		return "", errors.New("missing " + name + " parameter")
	}
	return val, nil
}
