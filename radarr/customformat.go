package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// CustomFormatService provides methods for the /customformat endpoint.
type CustomFormatService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CustomFormatBulkRequest is the body for bulk update and delete operations on custom formats.
type CustomFormatBulkRequest struct {
	IDs                             []int `json:"ids,omitempty"`
	IncludeCustomFormatWhenRenaming *bool `json:"includeCustomFormatWhenRenaming,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all custom formats.
func (s *CustomFormatService) List(ctx context.Context) ([]CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]CustomFormat{}).
		Get("/api/v3/customformat")
	if err != nil {
		return nil, fmt.Errorf("radarr: list custom formats: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list custom formats: %w", err)
	}

	return derefResult[[]CustomFormat](resp)
}

// Get returns a single custom format by ID.
func (s *CustomFormatService) Get(ctx context.Context, id int) (*CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFormat{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/customformat/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get custom format %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get custom format %d: %w", id, err)
	}

	return ptrResult[CustomFormat](resp)
}

// Schema returns the available specification schemas for custom formats.
func (s *CustomFormatService) Schema(ctx context.Context) ([]CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]CustomFormat{}).
		Get("/api/v3/customformat/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: custom format schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: custom format schema: %w", err)
	}

	return derefResult[[]CustomFormat](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new custom format.
func (s *CustomFormatService) Create(ctx context.Context, body CustomFormat) (*CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFormat{}).
		SetBody(body).
		Post("/api/v3/customformat")
	if err != nil {
		return nil, fmt.Errorf("radarr: create custom format: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create custom format: %w", err)
	}

	return ptrResult[CustomFormat](resp)
}

// Update replaces a custom format by ID.
func (s *CustomFormatService) Update(ctx context.Context, id int, body CustomFormat) (*CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CustomFormat{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/customformat/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update custom format %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update custom format %d: %w", id, err)
	}

	return ptrResult[CustomFormat](resp)
}

// UpdateBulk updates multiple custom formats at once.
func (s *CustomFormatService) UpdateBulk(ctx context.Context, body CustomFormatBulkRequest) ([]CustomFormat, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]CustomFormat{}).
		SetBody(body).
		Put("/api/v3/customformat/bulk")
	if err != nil {
		return nil, fmt.Errorf("radarr: update custom formats bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update custom formats bulk: %w", err)
	}

	return derefResult[[]CustomFormat](resp)
}

// Delete removes a custom format by ID.
func (s *CustomFormatService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/customformat/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete custom format %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete custom format %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple custom formats in a single request.
func (s *CustomFormatService) DeleteBulk(ctx context.Context, body CustomFormatBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/customformat/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete custom formats bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete custom formats bulk: %w", err)
	}

	return nil
}
