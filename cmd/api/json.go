package main

import (
	"encoding/json"
	"net/http"
)

// Write JSON func will set status(custom),return json encoder and return error
func writeJSON(w http.ResponseWriter, status int,data any) error{
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}