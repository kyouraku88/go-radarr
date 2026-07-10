package radarr

import (
	"context"
	"encoding/json"
	"fmt"
)

// LocalizationService provides methods for the /localization endpoint.
type LocalizationService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// LocalizationLanguage identifies the current UI localization language.
type LocalizationLanguage struct {
	Identifier *string `json:"identifier,omitempty"`
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

// Get returns the localization strings for the current UI language as a
// key→value map (e.g. "Add" → "Add").
func (s *LocalizationService) Get(ctx context.Context) (map[string]string, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		Get("/api/v3/localization")
	if err != nil {
		return nil, fmt.Errorf("radarr: get localization: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get localization: %w", err)
	}

	// The endpoint returns a flat JSON object of string→string.
	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("radarr: get localization: %w", err)
	}

	return result, nil
}

// GetLanguage returns the identifier of the active localization language.
func (s *LocalizationService) GetLanguage(ctx context.Context) (*LocalizationLanguage, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&LocalizationLanguage{}).
		Get("/api/v3/localization/language")
	if err != nil {
		return nil, fmt.Errorf("radarr: get localization language: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get localization language: %w", err)
	}

	return ptrResult[LocalizationLanguage](resp)
}
