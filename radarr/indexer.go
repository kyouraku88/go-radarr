package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// IndexerService provides methods for the /indexer endpoint.
type IndexerService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Indexer is a configured indexer (e.g. Jackett, Prowlarr).
type Indexer struct {
	ID                      int              `json:"id"`
	Name                    *string          `json:"name,omitempty"`
	Fields                  []Field          `json:"fields,omitempty"`
	ImplementationName      *string          `json:"implementationName,omitempty"`
	Implementation          *string          `json:"implementation,omitempty"`
	ConfigContract          *string          `json:"configContract,omitempty"`
	InfoLink                *string          `json:"infoLink,omitempty"`
	Message                 *ProviderMessage `json:"message,omitempty"`
	Tags                    []int            `json:"tags,omitempty"`
	Presets                 []Indexer        `json:"presets,omitempty"`
	EnableRss               bool             `json:"enableRss"`
	EnableAutomaticSearch   bool             `json:"enableAutomaticSearch"`
	EnableInteractiveSearch bool             `json:"enableInteractiveSearch"`
	SupportsRss             bool             `json:"supportsRss"`
	SupportsSearch          bool             `json:"supportsSearch"`
	Protocol                DownloadProtocol `json:"protocol,omitempty"`
	Priority                int              `json:"priority"`
	DownloadClientID        int              `json:"downloadClientId"`
}

// IndexerBulkRequest is the body for bulk update and delete operations on indexers.
type IndexerBulkRequest struct {
	IDs                     []int     `json:"ids,omitempty"`
	Tags                    []int     `json:"tags,omitempty"`
	ApplyTags               ApplyTags `json:"applyTags,omitempty"`
	EnableRss               *bool     `json:"enableRss,omitempty"`
	EnableAutomaticSearch   *bool     `json:"enableAutomaticSearch,omitempty"`
	EnableInteractiveSearch *bool     `json:"enableInteractiveSearch,omitempty"`
	Priority                *int      `json:"priority,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get / Schema
// ---------------------------------------------------------------------------

// List returns all configured indexers.
func (s *IndexerService) List(ctx context.Context) ([]Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Indexer{}).
		Get("/api/v3/indexer")
	if err != nil {
		return nil, fmt.Errorf("radarr: list indexers: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list indexers: %w", err)
	}

	return derefResult[[]Indexer](resp)
}

// Get returns a single indexer by ID.
func (s *IndexerService) Get(ctx context.Context, id int) (*Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Indexer{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/indexer/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get indexer %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get indexer %d: %w", id, err)
	}

	return ptrResult[Indexer](resp)
}

// Schema returns the available indexer implementation schemas.
func (s *IndexerService) Schema(ctx context.Context) ([]Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Indexer{}).
		Get("/api/v3/indexer/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: indexer schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: indexer schema: %w", err)
	}

	return derefResult[[]Indexer](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new indexer.
func (s *IndexerService) Create(ctx context.Context, body Indexer) (*Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Indexer{}).
		SetBody(body).
		Post("/api/v3/indexer")
	if err != nil {
		return nil, fmt.Errorf("radarr: create indexer: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create indexer: %w", err)
	}

	return ptrResult[Indexer](resp)
}

// Update replaces an indexer by ID.
func (s *IndexerService) Update(ctx context.Context, id int, body Indexer) (*Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Indexer{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/indexer/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update indexer %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update indexer %d: %w", id, err)
	}

	return ptrResult[Indexer](resp)
}

// UpdateBulk updates multiple indexers at once.
func (s *IndexerService) UpdateBulk(ctx context.Context, body IndexerBulkRequest) ([]Indexer, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Indexer{}).
		SetBody(body).
		Put("/api/v3/indexer/bulk")
	if err != nil {
		return nil, fmt.Errorf("radarr: update indexers bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update indexers bulk: %w", err)
	}

	return derefResult[[]Indexer](resp)
}

// Delete removes an indexer by ID.
func (s *IndexerService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/indexer/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete indexer %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete indexer %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple indexers in a single request.
func (s *IndexerService) DeleteBulk(ctx context.Context, body IndexerBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/indexer/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete indexers bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete indexers bulk: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Test / Action
// ---------------------------------------------------------------------------

// Test validates an indexer configuration.
func (s *IndexerService) Test(ctx context.Context, body Indexer) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/indexer/test")
	if err != nil {
		return fmt.Errorf("radarr: test indexer: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test indexer: %w", err)
	}

	return nil
}

// TestAll validates all configured indexers.
func (s *IndexerService) TestAll(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/indexer/testall")
	if err != nil {
		return fmt.Errorf("radarr: test all indexers: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test all indexers: %w", err)
	}

	return nil
}

// Action executes a named action on an indexer.
func (s *IndexerService) Action(ctx context.Context, name string, body Indexer) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("name", name).
		SetBody(body).
		Post("/api/v3/indexer/action/{name}")
	if err != nil {
		return fmt.Errorf("radarr: indexer action %s: %w", name, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: indexer action %s: %w", name, err)
	}

	return nil
}
