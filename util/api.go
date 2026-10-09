package util

import (
	"errors"
	"main/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError remains an alias for server application errors.
type APIError = api.APIError

func WriteAPIError(c *gin.Context, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		apiErr = &APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "Unable to complete request"}
	}
	c.JSON(apiErr.Status, api.ErrorResponse{Error: apiErr})
}

var (
	ErrInvalidRequest   = &APIError{Status: 400, Code: "invalid_request", Message: "Invalid request body"}
	ErrUnauthorized     = &APIError{Status: 401, Code: "unauthorized", Message: "Missing or invalid authentication token"}
	ErrForbidden        = &APIError{Status: 403, Code: "forbidden", Message: "You do not have permission for this operation"}
	ErrNotFound         = &APIError{Status: 404, Code: "not_found", Message: "Resource not found"}
	ErrLobbyClosed      = &APIError{Status: 409, Code: "lobby_closed", Message: "This lobby has already launched"}
	ErrUnsupportedState = &APIError{Status: 409, Code: "unsupported_state", Message: "The saved match schema, rules or catalogue version is not supported"}
	ErrStaleRevision    = &APIError{Status: 409, Code: "stale_revision", Message: "The match has changed; refresh its state before choosing"}
)
