package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func HandleErrors(h ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		log.Println("[ERROR]", err)

		code := http.StatusInternalServerError
		msg := "Internal Server Error"

		if httpErr, ok := err.(*HTTPError); ok {
			code = httpErr.Code
			msg = httpErr.Message
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
	}
}
