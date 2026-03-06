package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Success: true, Data: data})
}

func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Success: false, Error: msg})
}

func ValidationError(w http.ResponseWriter, errs any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(Response{Success: false, Error: "validation failed", Errors: errs})
}

func OK(w http.ResponseWriter, data any)              { JSON(w, http.StatusOK, data) }
func Created(w http.ResponseWriter, data any)         { JSON(w, http.StatusCreated, data) }
func NoContent(w http.ResponseWriter)                 { w.WriteHeader(http.StatusNoContent) }
func BadRequest(w http.ResponseWriter, msg string)    { Error(w, http.StatusBadRequest, msg) }
func Unauthorized(w http.ResponseWriter, msg string)  { Error(w, http.StatusUnauthorized, msg) }
func Forbidden(w http.ResponseWriter, msg string)     { Error(w, http.StatusForbidden, msg) }
func NotFound(w http.ResponseWriter, msg string)      { Error(w, http.StatusNotFound, msg) }
func Conflict(w http.ResponseWriter, msg string)      { Error(w, http.StatusConflict, msg) }
func TooManyRequests(w http.ResponseWriter, msg string) { Error(w, http.StatusTooManyRequests, msg) }
func InternalError(w http.ResponseWriter)             { Error(w, http.StatusInternalServerError, "internal server error") }
