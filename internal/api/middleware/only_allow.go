package middleware

import (
	"net/http"
)

func OnlyAllow(method string, next ErrorHandlerFunc) ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != method {
			return NewHttpError(http.StatusMethodNotAllowed, "Method not allowed")
		}
		return next(w, r)
	}
}
