package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf(
		"Internal server error: %s path: %s error: %s",
		r.Method,
		r.URL.Path,
		err,
	)
	jsonErrorResponse(w, http.StatusInternalServerError, "The server encountered a problem.")
}

func (app *application) resourceConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf(
		"Data conflict: %s path: %s error: %s",
		r.Method,
		r.URL.Path,
		err,
	)
	jsonErrorResponse(w, http.StatusConflict, "There is data conflict. Try again.")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf(
		"Bad request: %s path: %s error: %s",
		r.Method,
		r.URL.Path,
		err,
	)
	jsonErrorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) resourceNotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf(
		"Resource not found: %s path: %s error: %s",
		r.Method,
		r.URL.Path,
		err,
	)
	jsonErrorResponse(w, http.StatusNotFound, "Not found")
}
