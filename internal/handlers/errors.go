package handlers

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func errorResponse(w http.ResponseWriter, status int, message any) {
	respEnvelope := envelope{
		"error": message,
	}

	if err := writeJSON(w, status, respEnvelope, nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func serverErrorResponse(w http.ResponseWriter, err error) {
	log.Println(err)
	message := "a problem has occured on the server, try again later"
	errorResponse(w, http.StatusInternalServerError, message)
}

func notFoundResponse(w http.ResponseWriter) {
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func badRequestResponse(w http.ResponseWriter, err error) {
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func invalidCredentialsResponse(w http.ResponseWriter) {
	message := "invalid credentials"
	errorResponse(w, http.StatusUnauthorized, message)
}

func notPermittedResponse(w http.ResponseWriter) {
	message := "you do not have permission to access this resource"
	errorResponse(w, http.StatusForbidden, message)
}

func failedValidationResponse(w http.ResponseWriter, err error) {
	errs := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = err.Tag()
	}

	message := "the request was invalid"
	respEvenlope := envelope{
		"error":  message,
		"fields": errs,
	}
	errorResponse(w, http.StatusUnprocessableEntity, respEvenlope)
}
