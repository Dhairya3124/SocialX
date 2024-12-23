package main

import (
	"encoding/json"
	"net/http"
)

// Write JSON func will set status(custom),return jsonEncoder
func writeJSON(w http.ResponseWriter, status int,data any) error{
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}
// ReadJSON will read the client json and setup the max required storage needed to scan
func readJSON(w http.ResponseWriter,r *http.Request,data any) error{
	maxBytes:=1_048_578
	r.Body = http.MaxBytesReader(w,r.Body,int64(maxBytes))
	decoder:=json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}
//writeJSONError: will return the error json
func writeJSONError(w http.ResponseWriter,status int, message string)error{
	type envelope struct{
		Error any `json:"error"`
	}
	return writeJSON(w,status,&envelope{Error: message})
}
// jsonResponse function is simple func which returns in the format:  data:
func jsonResponse(w http.ResponseWriter, status int,data any) error{
	type envelope struct{
		Data any `json:"data"`
	}
	return writeJSON(w,status,&envelope{Data:data})

}