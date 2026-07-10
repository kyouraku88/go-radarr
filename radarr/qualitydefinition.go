package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// QualityDefinitionService provides methods for the /qualitydefinition endpoint.
type QualityDefinitionService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// QualityDefinition defines the size limits for a quality tier.
type QualityDefinition struct {
	ID            int      `json:"id"`
	Quality       *Quality `json:"quality,omitempty"`
	Title         *string  `json:"title,omitempty"`
	Weight        int      `json:"weight"`
	MinSize       *float64 `json:"minSize,omitempty"`
	MaxSize       *float64 `json:"maxSize,omitempty"`
	PreferredSize *float64 `json:"preferredSize,omitempty"`
}

// QualityDefinitionLimits is the absolute min/max sizes allowed across all definitions.
type QualityDefinitionLimits struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all quality definitions.
func (s *QualityDefinitionService) List(ctx context.Context) ([]QualityDefinition, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]QualityDefinition{}).
		Get("/api/v3/qualitydefinition")
	if err != nil {
		return nil, fmt.Errorf("radarr: list quality definitions: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list quality definitions: %w", err)
	}

	return derefResult[[]QualityDefinition](resp)
}

// Get returns a single quality definition by ID.
func (s *QualityDefinitionService) Get(ctx context.Context, id int) (*QualityDefinition, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityDefinition{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/qualitydefinition/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get quality definition %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get quality definition %d: %w", id, err)
	}

	return ptrResult[QualityDefinition](resp)
}

// Limits returns the allowed min/max size range for quality definitions.
func (s *QualityDefinitionService) Limits(ctx context.Context) (*QualityDefinitionLimits, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityDefinitionLimits{}).
		Get("/api/v3/qualitydefinition/limits")
	if err != nil {
		return nil, fmt.Errorf("radarr: get quality definition limits: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get quality definition limits: %w", err)
	}

	return ptrResult[QualityDefinitionLimits](resp)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// Update replaces a single quality definition by ID.
func (s *QualityDefinitionService) Update(ctx context.Context, id int, body QualityDefinition) (*QualityDefinition, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityDefinition{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/qualitydefinition/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update quality definition %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update quality definition %d: %w", id, err)
	}

	return ptrResult[QualityDefinition](resp)
}

// UpdateBulk replaces all quality definitions at once.
func (s *QualityDefinitionService) UpdateBulk(ctx context.Context, body []QualityDefinition) ([]QualityDefinition, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]QualityDefinition{}).
		SetBody(body).
		Put("/api/v3/qualitydefinition/update")
	if err != nil {
		return nil, fmt.Errorf("radarr: update quality definitions bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update quality definitions bulk: %w", err)
	}

	return derefResult[[]QualityDefinition](resp)
}
