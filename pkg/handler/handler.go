package handler

import (
	"net/http"
)

// Ping handler for returning health status
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
