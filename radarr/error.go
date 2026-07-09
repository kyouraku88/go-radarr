package radarr

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// APIError is returned when the Radarr API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Status     string
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("radarr: HTTP %d: %s", e.StatusCode, e.Message)
	}

	return "radarr: HTTP " + e.Status
}

func checkResponse(resp *resty.Response) error {
	if resp.IsSuccess() {
		return nil
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
	}

	var body []struct {
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal(resp.Body(), &body); err == nil && len(body) > 0 {
		apiErr.Message = body[0].ErrorMessage
	}

	return apiErr
}
