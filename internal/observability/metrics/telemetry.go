package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Init sets up default metrics, registers runtime & process collectors, etc.
// Returns the HTTP handler to mount at /metrics.
func Init() http.Handler {
	// You can create your own registry if you don't want global default.
	// For simplicity, using the default.

	// For registering custom metrics refer to
	// https://prometheus.io/docs/guides/go-application/

	// Return the HTTP handler
	return promhttp.Handler()
}
