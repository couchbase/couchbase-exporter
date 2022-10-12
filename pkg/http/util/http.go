package util

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/couchbase/couchbase-exporter/pkg/log"
)

// Respond generate a response for a http request.
func Respond(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		RespondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Debug("error buffer writing success response: %v", err)
	}
}

// RespondErr generate a response error for a http response.
func RespondErr(w http.ResponseWriter, r *http.Request, err error, code int) {
	errObj := struct {
		Error string `json:"error"`
	}{Error: err.Error()}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		log.Debug("error encoding error response: %v", err)
	}
}
