package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorMessage structure that returns errors
type ErrorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Map is a convenient way to create objects of unknown types.
type Map map[string]interface{}

// JSON standardized JSON response.
func JSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) error {
	if data == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		return nil
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(j)
	return nil
}

// HTTPError standardized error response in JSON format.
func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) error {
	msg := ErrorMessage{
		Status:  statusCode,
		Message: message,
	}

	return JSON(w, r, statusCode, msg)
}

// HTTPErrors standardized errors response in JSON format.
func HTTPErrors(w http.ResponseWriter, r *http.Request, statusCode int, errors interface{}) error {

	marshalErrors, err := json.Marshal(errors)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(marshalErrors)

	return JSON(w, r, statusCode, errors)
}
