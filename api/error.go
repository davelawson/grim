package api

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
