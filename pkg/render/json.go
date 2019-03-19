package render

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(logger *zap.Logger, w http.ResponseWriter, r *http.Request, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := enc.Encode(v); err != nil {
		logger.Error(
			"Error marshalling json",
			zap.Error(err),
			zapdriver.HTTP(zapdriver.NewHTTP(r, nil)),
		)

		w.WriteHeader(http.StatusInternalServerError)
		errorCoded, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		_, _ = w.Write(errorCoded)
		return
	}

	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}
