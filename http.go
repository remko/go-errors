package errors

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type apiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func FailJSON(w http.ResponseWriter, r *http.Request, err error) {
	c, m := Code(err)
	code, message := HTTPCodeMessage(c, m)
	payload := Payload(err)
	if payload == nil {
		payload = &apiError{Code: code, Message: message}
	}
	// log.Printf(r.Context(), lg.HTTPCodeSeverity(code), "error: %+v", err)
	writeJSON(w, code, payload)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}
