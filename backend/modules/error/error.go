package error

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	IsClientError bool
	Errors        []string
	Writer        http.ResponseWriter
}

type ErrorResponse struct {
	IsClientError bool     `json:"isClientError"`
	Errors        []string `json:"errors"`
}

func NewError(IsClientError bool, w http.ResponseWriter) Error {
	return Error{IsClientError: IsClientError, Errors: make([]string, 0, 10), Writer: w}
}

func (e *Error) NewError(error string) {
	e.Errors = append(e.Errors, error)
}

func (e *Error) ErrorsExist() bool {
	return len(e.Errors) > 0
}

func (e *Error) ThrowInternalError() {
	e.IsClientError = false
	e.Errors = make([]string, 0)
	e.ThrowError()
}

func (e *Error) ThrowError() {
	if e.IsClientError {
		e.Writer.WriteHeader(400)
	} else {
		e.Writer.WriteHeader(500)
	}

	response := ErrorResponse{
		IsClientError: e.IsClientError,
		Errors:        e.Errors,
	}

	json.NewEncoder(e.Writer).Encode(response)
}
