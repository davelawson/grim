package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError is an application outcome safe to expose to a client.
type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

type ErrorResponse struct {
	Error *APIError `json:"error"`
}

func WriteAPIError(c *gin.Context, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		apiErr = &APIError{http.StatusInternalServerError, "internal_error", "Unable to complete request"}
	}
	c.JSON(apiErr.Status, ErrorResponse{Error: apiErr})
}

var (
	ErrInvalidRequest   = &APIError{400, "invalid_request", "Invalid request body"}
	ErrUnauthorized     = &APIError{401, "unauthorized", "Missing or invalid authentication token"}
	ErrForbidden        = &APIError{403, "forbidden", "You do not have permission for this operation"}
	ErrNotFound         = &APIError{404, "not_found", "Resource not found"}
	ErrLobbyClosed      = &APIError{409, "lobby_closed", "This lobby has already launched"}
	ErrUnsupportedState = &APIError{409, "unsupported_state", "The saved match schema, rules or catalogue version is not supported"}
	ErrStaleRevision    = &APIError{409, "stale_revision", "The match has changed; refresh its state before choosing"}
)
