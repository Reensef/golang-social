package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func jsonErrorResponse(w http.ResponseWriter, status int, message string) error {
	return writeJson(w, status, map[string]string{"error": message})
}


func jsonDataResponse(w http.ResponseWriter, status int, data any) error {
	type envelop struct {
		Data any `json:"data"`
	}
	return writeJson(w, status, &envelop{Data: data})
}
