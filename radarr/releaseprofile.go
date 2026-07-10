package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// ReleaseProfileService provides methods for the /releaseprofile endpoint.
type ReleaseProfileService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ReleaseProfile defines required and ignored terms for release selection.
type ReleaseProfile struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	Enabled   bool    `json:"enabled"`
	Required  any     `json:"required,omitempty"`
	Ignored   any     `json:"ignored,omitempty"`
	IndexerID int     `json:"indexerId"`
	Tags      []int   `json:"tags,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all release profiles.
func (s *ReleaseProfileService) List(ctx context.Context) ([]ReleaseProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ReleaseProfile{}).
		Get("/api/v3/releaseprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: list release profiles: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list release profiles: %w", err)
	}

	return derefResult[[]ReleaseProfile](resp)
}

// Get returns a single release profile by ID.
func (s *ReleaseProfileService) Get(ctx context.Context, id int) (*ReleaseProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ReleaseProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/releaseprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get release profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get release profile %d: %w", id, err)
	}

	return ptrResult[ReleaseProfile](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new release profile.
func (s *ReleaseProfileService) Create(ctx context.Context, body ReleaseProfile) (*ReleaseProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ReleaseProfile{}).
		SetBody(body).
		Post("/api/v3/releaseprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: create release profile: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create release profile: %w", err)
	}

	return ptrResult[ReleaseProfile](resp)
}

// Update replaces a release profile by ID.
func (s *ReleaseProfileService) Update(ctx context.Context, id int, body ReleaseProfile) (*ReleaseProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ReleaseProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/releaseprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update release profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update release profile %d: %w", id, err)
	}

	return ptrResult[ReleaseProfile](resp)
}

// Delete removes a release profile by ID.
func (s *ReleaseProfileService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/releaseprofile/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete release profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete release profile %d: %w", id, err)
	}

	return nil
}
