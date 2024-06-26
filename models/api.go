package models

import (
	"net/http"

	"github.com/google/uuid"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

type UpdateUserStruct struct {
	ID    uuid.UUID `json:"id"`
	Score int       `json:"score"`
}
