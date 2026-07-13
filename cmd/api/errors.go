package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path:%s error: %s", r.Method, r.URL.Path, err.Error())

	writeJsonError(w, http.StatusInternalServerError, "Internal server error problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s path:%s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusBadRequest, "Bad request error problem")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized error: %s path:%s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusUnauthorized, "Unauthorized error problem")
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s path:%s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusNotFound, "Status not found error problem")
}

type ErrorResponse struct {
    Error string `json:"error"`
}
func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, message string) {
	err := writeJSON(w, http.StatusConflict, ErrorResponse{
    Error: message,
})

	if err != nil {
		app.internalServerError(w, r, err)
	}
}
