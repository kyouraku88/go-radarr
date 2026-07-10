package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// LanguageService provides methods for the /language endpoint.
type LanguageService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// LanguageResource is the full language resource returned by /language.
type LanguageResource struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	NameLower *string `json:"nameLower,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all available languages.
func (s *LanguageService) List(ctx context.Context) ([]LanguageResource, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]LanguageResource{}).
		Get("/api/v3/language")
	if err != nil {
		return nil, fmt.Errorf("radarr: list languages: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list languages: %w", err)
	}

	return derefResult[[]LanguageResource](resp)
}

// Get returns a single language by ID.
func (s *LanguageService) Get(ctx context.Context, id int) (*LanguageResource, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&LanguageResource{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/language/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get language %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get language %d: %w", id, err)
	}

	return ptrResult[LanguageResource](resp)
}
