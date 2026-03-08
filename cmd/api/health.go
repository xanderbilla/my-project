package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request)  {
	data := map[string]string{
		"status": "ok",
		"env": app.config.env,
		"version": version,
	}

	if err := Success(w, http.StatusOK, "OK", data); err != nil {
		Error(w, http.StatusInternalServerError, "Something went wrong")
	}
}