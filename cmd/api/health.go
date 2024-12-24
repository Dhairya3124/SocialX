package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err:=app.jsonResponse(w,http.StatusOK,"OK");err!=nil{
		writeJSONError(w,http.StatusInternalServerError,err.Error())
	}
}
