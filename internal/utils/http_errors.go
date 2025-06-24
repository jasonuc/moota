package utils

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ErrorResponse(w http.ResponseWriter, status int, message any) {
	respEnvelope := Envelope{
		"error": message,
	}

	if err := WriteJSON(w, status, respEnvelope, nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func ServerErrorResponse(w http.ResponseWriter, err error) {
	log.Println(err)
	message := "A problem has occured on the server, try again later"
	ErrorResponse(w, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter) {
	message := "Requested resource could not be found"
	ErrorResponse(w, http.StatusNotFound, message)
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusBadRequest, err.Error())
}

func InvalidCredentialsResponse(w http.ResponseWriter) {
	message := "Invalid credentials"
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func NotPermittedResponse(w http.ResponseWriter) {
	message := "You do not have permission to access this resource"
	ErrorResponse(w, http.StatusForbidden, message)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	message := "Unauthorized"
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func FailedValidationResponse(w http.ResponseWriter, err error) {
	errs := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = err.Tag()
	}

	message := "Request was invalid"
	respEvenlope := Envelope{
		"message": message,
		"fields":  errs,
	}
	ErrorResponse(w, http.StatusUnprocessableEntity, respEvenlope)
}
