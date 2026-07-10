package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// CustomFilterService provides methods for the /customfilter endpoint.
type CustomFilterService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CustomFilter is a user-defined saved filter for a UI view.
type CustomFilter struct {
	ID      int              `json:"id"`
	Type    *string          `json:"type,omitempty"`
	Label   *string          `json:"label,omitempty"`
	Filters []map[string]any `json:"filters,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all custom filters.
func (s *CustomFilterService) List(ctx context.Context) ([]CustomFilter, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]CustomFilter{}).
		Get("/api/v3/customfilter")
	if err != nil {
		return nil, fmt.Errorf("radarr: list custom filters: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list custom filters: %w", err)
	}

	return derefResult[[]CustomFilter](resp)
}

// Get returns a single custom filter by ID.
func (s *CustomFilterService) Get(ctx context.Context, id int) (*CustomFilter, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFilter{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/customfilter/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get custom filter %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get custom filter %d: %w", id, err)
	}

	return ptrResult[CustomFilter](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new custom filter.
func (s *CustomFilterService) Create(ctx context.Context, body CustomFilter) (*CustomFilter, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFilter{}).
		SetBody(body).
		Post("/api/v3/customfilter")
	if err != nil {
		return nil, fmt.Errorf("radarr: create custom filter: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create custom filter: %w", err)
	}

	return ptrResult[CustomFilter](resp)
}

// Update replaces a custom filter by ID.
func (s *CustomFilterService) Update(ctx context.Context, id int, body CustomFilter) (*CustomFilter, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFilter{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/customfilter/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update custom filter %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update custom filter %d: %w", id, err)
	}

	return ptrResult[CustomFilter](resp)
}

// Delete removes a custom filter by ID.
func (s *CustomFilterService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/customfilter/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete custom filter %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete custom filter %d: %w", id, err)
	}

	return nil
}
