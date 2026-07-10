package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// DownloadClientService provides methods for the /downloadclient endpoint.
type DownloadClientService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// DownloadClient is a configured download client (e.g. qBittorrent, SABnzbd).
type DownloadClient struct {
	ID                       int              `json:"id"`
	Name                     *string          `json:"name,omitempty"`
	Fields                   []Field          `json:"fields,omitempty"`
	ImplementationName       *string          `json:"implementationName,omitempty"`
	Implementation           *string          `json:"implementation,omitempty"`
	ConfigContract           *string          `json:"configContract,omitempty"`
	InfoLink                 *string          `json:"infoLink,omitempty"`
	Message                  *ProviderMessage `json:"message,omitempty"`
	Tags                     []int            `json:"tags,omitempty"`
	Presets                  []DownloadClient `json:"presets,omitempty"`
	Enable                   bool             `json:"enable"`
	Protocol                 DownloadProtocol `json:"protocol,omitempty"`
	Priority                 int              `json:"priority"`
	RemoveCompletedDownloads bool             `json:"removeCompletedDownloads"`
	RemoveFailedDownloads    bool             `json:"removeFailedDownloads"`
}

// DownloadClientBulkRequest is the body for bulk update and delete operations.
type DownloadClientBulkRequest struct {
	IDs                      []int     `json:"ids,omitempty"`
	Tags                     []int     `json:"tags,omitempty"`
	ApplyTags                ApplyTags `json:"applyTags,omitempty"`
	Enable                   *bool     `json:"enable,omitempty"`
	Priority                 *int      `json:"priority,omitempty"`
	RemoveCompletedDownloads *bool     `json:"removeCompletedDownloads,omitempty"`
	RemoveFailedDownloads    *bool     `json:"removeFailedDownloads,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get / Schema
// ---------------------------------------------------------------------------

// List returns all configured download clients.
func (s *DownloadClientService) List(ctx context.Context) ([]DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]DownloadClient{}).
		Get("/api/v3/downloadclient")
	if err != nil {
		return nil, fmt.Errorf("radarr: list download clients: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list download clients: %w", err)
	}

	return derefResult[[]DownloadClient](resp)
}

// Get returns a single download client by ID.
func (s *DownloadClientService) Get(ctx context.Context, id int) (*DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DownloadClient{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/downloadclient/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get download client %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get download client %d: %w", id, err)
	}

	return ptrResult[DownloadClient](resp)
}

// Schema returns the available download client implementation schemas.
func (s *DownloadClientService) Schema(ctx context.Context) ([]DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]DownloadClient{}).
		Get("/api/v3/downloadclient/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: download client schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: download client schema: %w", err)
	}

	return derefResult[[]DownloadClient](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new download client.
func (s *DownloadClientService) Create(ctx context.Context, body DownloadClient) (*DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DownloadClient{}).
		SetBody(body).
		Post("/api/v3/downloadclient")
	if err != nil {
		return nil, fmt.Errorf("radarr: create download client: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create download client: %w", err)
	}

	return ptrResult[DownloadClient](resp)
}

// Update replaces a download client by ID.
func (s *DownloadClientService) Update(ctx context.Context, id int, body DownloadClient) (*DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DownloadClient{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/downloadclient/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update download client %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update download client %d: %w", id, err)
	}

	return ptrResult[DownloadClient](resp)
}

// UpdateBulk updates multiple download clients at once.
func (s *DownloadClientService) UpdateBulk(ctx context.Context, body DownloadClientBulkRequest) ([]DownloadClient, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]DownloadClient{}).
		SetBody(body).
		Put("/api/v3/downloadclient/bulk")
	if err != nil {
		return nil, fmt.Errorf("radarr: update download clients bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update download clients bulk: %w", err)
	}

	return derefResult[[]DownloadClient](resp)
}

// Delete removes a download client by ID.
func (s *DownloadClientService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/downloadclient/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete download client %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete download client %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple download clients in a single request.
func (s *DownloadClientService) DeleteBulk(ctx context.Context, body DownloadClientBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/downloadclient/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete download clients bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete download clients bulk: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Test / Action
// ---------------------------------------------------------------------------

// Test validates a download client configuration.
func (s *DownloadClientService) Test(ctx context.Context, body DownloadClient) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/downloadclient/test")
	if err != nil {
		return fmt.Errorf("radarr: test download client: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test download client: %w", err)
	}

	return nil
}

// TestAll validates all configured download clients.
func (s *DownloadClientService) TestAll(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/downloadclient/testall")
	if err != nil {
		return fmt.Errorf("radarr: test all download clients: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test all download clients: %w", err)
	}

	return nil
}

// Action executes a named action on a download client (e.g. "checkHealth").
func (s *DownloadClientService) Action(ctx context.Context, name string, body DownloadClient) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("name", name).
		SetBody(body).
		Post("/api/v3/downloadclient/action/{name}")
	if err != nil {
		return fmt.Errorf("radarr: download client action %s: %w", name, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: download client action %s: %w", name, err)
	}

	return nil
}
