package responses

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func OK(w http.ResponseWriter, response interface{}) {
	answer(w, http.StatusOK, response)
}

func Created(w http.ResponseWriter, response interface{}) {
	answer(w, http.StatusCreated, response)
}

func Duplicated(w http.ResponseWriter, response string) {
	Abort(w, http.StatusConflict, response)
}

func BadRequest(w http.ResponseWriter, response string) {
	Abort(w, http.StatusBadRequest, response)
}

func NotFound(w http.ResponseWriter, response string) {
	Abort(w, http.StatusNotFound, response)
}

func Error(w http.ResponseWriter, err error) {
	if e, ok := errors.Cause(err).(interface {
		StatusCode() int
	}); ok {
		answer(w, e.StatusCode(), e)
		return
	}
	Abort(w, http.StatusInternalServerError, err.Error())
}

func Abort(w http.ResponseWriter, status int, message string) {
	answer(w, status, map[string]interface{}{
		"status":  status,
		"error":   http.StatusText(status),
		"message": message,
		"cause":   make([]string, 0),
	})
}

func answer(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

