package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Envelope struct {
	Data    any     `json:"data"`
	Success bool    `json:"success"`
	Status  int     `json:"status"`
	Errors  []Error `json:"errors"`
}

func Success(w http.ResponseWriter, status int, data any) error {
	return write(w, status, Envelope{
		Data:    data,
		Success: true,
		Status:  status,
		Errors:  make([]Error, 0),
	})
}

func Failure(w http.ResponseWriter, status int, errors ...Error) error {
	if len(errors) == 0 {
		errors = []Error{GeneralError(http.StatusText(status))}
	}

	return write(w, status, Envelope{
		Data:    nil,
		Success: false,
		Status:  status,
		Errors:  errors,
	})
}

func FieldError(field, message string) Error {
	return Error{Field: field, Message: message}
}

func GeneralError(message string) Error {
	return Error{Message: message}
}

func write(w http.ResponseWriter, status int, envelope Envelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(envelope)
}
