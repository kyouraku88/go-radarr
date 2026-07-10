package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// AutoTaggingService provides methods for the /autotagging endpoint.
type AutoTaggingService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// AutoTaggingSpecificationSchema is a condition rule within an auto-tagging profile.
type AutoTaggingSpecificationSchema struct {
	ID                 int     `json:"id"`
	Name               *string `json:"name,omitempty"`
	Implementation     *string `json:"implementation,omitempty"`
	ImplementationName *string `json:"implementationName,omitempty"`
	Negate             bool    `json:"negate"`
	Required           bool    `json:"required"`
	Fields             []Field `json:"fields,omitempty"`
}

// AutoTagging is an auto-tagging profile that applies tags based on conditions.
type AutoTagging struct {
	ID                      int                              `json:"id"`
	Name                    *string                          `json:"name,omitempty"`
	RemoveTagsAutomatically bool                             `json:"removeTagsAutomatically"`
	Tags                    []int                            `json:"tags,omitempty"`
	Specifications          []AutoTaggingSpecificationSchema `json:"specifications,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all auto-tagging profiles.
func (s *AutoTaggingService) List(ctx context.Context) ([]AutoTagging, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]AutoTagging{}).
		Get("/api/v3/autotagging")
	if err != nil {
		return nil, fmt.Errorf("radarr: list auto tagging: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list auto tagging: %w", err)
	}

	return derefResult[[]AutoTagging](resp)
}

// Get returns a single auto-tagging profile by ID.
func (s *AutoTaggingService) Get(ctx context.Context, id int) (*AutoTagging, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&AutoTagging{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/autotagging/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get auto tagging %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get auto tagging %d: %w", id, err)
	}

	return ptrResult[AutoTagging](resp)
}

// Schema returns the available specification schemas for auto-tagging.
func (s *AutoTaggingService) Schema(ctx context.Context) ([]AutoTagging, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]AutoTagging{}).
		Get("/api/v3/autotagging/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: auto tagging schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: auto tagging schema: %w", err)
	}

	return derefResult[[]AutoTagging](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new auto-tagging profile.
func (s *AutoTaggingService) Create(ctx context.Context, body AutoTagging) (*AutoTagging, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&AutoTagging{}).
		SetBody(body).
		Post("/api/v3/autotagging")
	if err != nil {
		return nil, fmt.Errorf("radarr: create auto tagging: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create auto tagging: %w", err)
	}

	return ptrResult[AutoTagging](resp)
}

// Update replaces an auto-tagging profile by ID.
func (s *AutoTaggingService) Update(ctx context.Context, id int, body AutoTagging) (*AutoTagging, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&AutoTagging{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/autotagging/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update auto tagging %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update auto tagging %d: %w", id, err)
	}

	return ptrResult[AutoTagging](resp)
}

// Delete removes an auto-tagging profile by ID.
func (s *AutoTaggingService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/autotagging/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete auto tagging %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete auto tagging %d: %w", id, err)
	}

	return nil
}
