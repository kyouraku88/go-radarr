package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// MetadataService provides methods for the /metadata endpoint.
type MetadataService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// MetadataProvider is a configured metadata agent.
type MetadataProvider struct {
	ID                 int                `json:"id"`
	Name               *string            `json:"name,omitempty"`
	Fields             []Field            `json:"fields,omitempty"`
	ImplementationName *string            `json:"implementationName,omitempty"`
	Implementation     *string            `json:"implementation,omitempty"`
	ConfigContract     *string            `json:"configContract,omitempty"`
	InfoLink           *string            `json:"infoLink,omitempty"`
	Message            *ProviderMessage   `json:"message,omitempty"`
	Tags               []int              `json:"tags,omitempty"`
	Presets            []MetadataProvider `json:"presets,omitempty"`
	Enable             bool               `json:"enable"`
}

// ---------------------------------------------------------------------------
// List / Get / Schema
// ---------------------------------------------------------------------------

// List returns all configured metadata agents.
func (s *MetadataService) List(ctx context.Context) ([]MetadataProvider, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]MetadataProvider{}).
		Get("/api/v3/metadata")
	if err != nil {
		return nil, fmt.Errorf("radarr: list metadata: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list metadata: %w", err)
	}

	return derefResult[[]MetadataProvider](resp)
}

// Get returns a single metadata agent by ID.
func (s *MetadataService) Get(ctx context.Context, id int) (*MetadataProvider, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MetadataProvider{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/metadata/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get metadata %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get metadata %d: %w", id, err)
	}

	return ptrResult[MetadataProvider](resp)
}

// Schema returns the available metadata agent implementation schemas.
func (s *MetadataService) Schema(ctx context.Context) ([]MetadataProvider, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]MetadataProvider{}).
		Get("/api/v3/metadata/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: metadata schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: metadata schema: %w", err)
	}

	return derefResult[[]MetadataProvider](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new metadata agent.
func (s *MetadataService) Create(ctx context.Context, body MetadataProvider) (*MetadataProvider, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MetadataProvider{}).
		SetBody(body).
		Post("/api/v3/metadata")
	if err != nil {
		return nil, fmt.Errorf("radarr: create metadata: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create metadata: %w", err)
	}

	return ptrResult[MetadataProvider](resp)
}

// Update replaces a metadata agent by ID.
func (s *MetadataService) Update(ctx context.Context, id int, body MetadataProvider) (*MetadataProvider, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MetadataProvider{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/metadata/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update metadata %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update metadata %d: %w", id, err)
	}

	return ptrResult[MetadataProvider](resp)
}

// Delete removes a metadata agent by ID.
func (s *MetadataService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/metadata/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete metadata %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete metadata %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Test / Action
// ---------------------------------------------------------------------------

// Test validates a metadata agent configuration.
func (s *MetadataService) Test(ctx context.Context, body MetadataProvider) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/metadata/test")
	if err != nil {
		return fmt.Errorf("radarr: test metadata: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test metadata: %w", err)
	}

	return nil
}

// TestAll validates all configured metadata agents.
func (s *MetadataService) TestAll(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/metadata/testall")
	if err != nil {
		return fmt.Errorf("radarr: test all metadata: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test all metadata: %w", err)
	}

	return nil
}

// Action executes a named action on a metadata agent.
func (s *MetadataService) Action(ctx context.Context, name string, body MetadataProvider) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("name", name).
		SetBody(body).
		Post("/api/v3/metadata/action/{name}")
	if err != nil {
		return fmt.Errorf("radarr: metadata action %s: %w", name, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: metadata action %s: %w", name, err)
	}

	return nil
}
