package maradocs

import "fmt"

// ApiError is the structured error returned by the API (JSON field api_error).
type ApiError struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

// HttpErrorResponse is the JSON body for non-2xx API errors.
type HttpErrorResponse struct {
	StatusCode int      `json:"status_code"`
	ApiErr     ApiError `json:"api_error"`
}

// APIException is returned when the API responds with a structured HTTP error body.
type APIException struct {
	StatusCode int
	Details    ApiError
}

func (e *APIException) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Details.Code, e.Details.Name, e.Details.Message)
}

// ValidationError is returned when validation-style endpoints report a logical error.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return "validation failed: " + e.Message
}

// ValidationVirus is returned when a virus is detected during validation.
type ValidationVirus struct {
	Virus string
}

func (e *ValidationVirus) Error() string {
	return "virus detected: " + e.Virus
}
