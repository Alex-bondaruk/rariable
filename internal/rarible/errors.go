package rarible

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError is a non-2xx response from the Rarible API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("rarible: %d %s", e.StatusCode, e.Code)
	}
	return fmt.Sprintf("rarible: %d %s: %s", e.StatusCode, e.Code, e.Message)
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func apiErrorFrom(resp *http.Response) *APIError {
	e := &APIError{
		StatusCode: resp.StatusCode,
		Code:       http.StatusText(resp.StatusCode),
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return e
	}

	var body errorBody
	if err := json.Unmarshal(raw, &body); err == nil && body.Code != "" {
		e.Code = body.Code
		e.Message = body.Message
		return e
	}

	e.Message = strings.TrimSpace(string(raw))
	return e
}
