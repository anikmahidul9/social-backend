package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJSON(w, http.StatusOK, map[string]string{"status": "available"});
	 err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
}
