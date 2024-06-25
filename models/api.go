package models

import "net/http"

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}
