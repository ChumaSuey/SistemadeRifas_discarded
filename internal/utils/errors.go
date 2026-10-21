package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrInternalServer      = errors.New("internal server error")
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message, detail string) *APIError {
	return &APIError{Code: code, Message: message, Detail: detail}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		WriteJSON(w, apiErr.Code, apiErr)
		return
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		WriteJSON(w, http.StatusBadRequest, NewAPIError(http.StatusBadRequest, ErrBadRequest.Error(), err.Error()))
	case errors.Is(err, ErrUnauthorized):
		WriteJSON(w, http.StatusUnauthorized, NewAPIError(http.StatusUnauthorized, ErrUnauthorized.Error(), err.Error()))
	case errors.Is(err, ErrForbidden):
		WriteJSON(w, http.StatusForbidden, NewAPIError(http.StatusForbidden, ErrForbidden.Error(), err.Error()))
	case errors.Is(err, ErrNotFound):
		WriteJSON(w, http.StatusNotFound, NewAPIError(http.StatusNotFound, ErrNotFound.Error(), err.Error()))
	case errors.Is(err, ErrConflict):
		WriteJSON(w, http.StatusConflict, NewAPIError(http.StatusConflict, ErrConflict.Error(), err.Error()))
	case errors.Is(err, ErrUnprocessableEntity):
		WriteJSON(w, http.StatusUnprocessableEntity, NewAPIError(http.StatusUnprocessableEntity, ErrUnprocessableEntity.Error(), err.Error()))
	default:
		WriteJSON(w, http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "internal server error", ""))
	}
}
