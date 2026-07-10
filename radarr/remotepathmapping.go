package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// RemotePathMappingService provides methods for the /remotepathmapping endpoint.
type RemotePathMappingService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// RemotePathMapping maps a remote download-client path to a local path.
type RemotePathMapping struct {
	ID         int     `json:"id"`
	Host       *string `json:"host,omitempty"`
	RemotePath *string `json:"remotePath,omitempty"`
	LocalPath  *string `json:"localPath,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all remote path mappings.
func (s *RemotePathMappingService) List(ctx context.Context) ([]RemotePathMapping, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]RemotePathMapping{}).
		Get("/api/v3/remotepathmapping")
	if err != nil {
		return nil, fmt.Errorf("radarr: list remote path mappings: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list remote path mappings: %w", err)
	}

	return derefResult[[]RemotePathMapping](resp)
}

// Get returns a single remote path mapping by ID.
func (s *RemotePathMappingService) Get(ctx context.Context, id int) (*RemotePathMapping, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&RemotePathMapping{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/remotepathmapping/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get remote path mapping %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get remote path mapping %d: %w", id, err)
	}

	return ptrResult[RemotePathMapping](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new remote path mapping.
func (s *RemotePathMappingService) Create(ctx context.Context, body RemotePathMapping) (*RemotePathMapping, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&RemotePathMapping{}).
		SetBody(body).
		Post("/api/v3/remotepathmapping")
	if err != nil {
		return nil, fmt.Errorf("radarr: create remote path mapping: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create remote path mapping: %w", err)
	}

	return ptrResult[RemotePathMapping](resp)
}

// Update replaces a remote path mapping by ID.
func (s *RemotePathMappingService) Update(ctx context.Context, id int, body RemotePathMapping) (*RemotePathMapping, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&RemotePathMapping{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/remotepathmapping/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update remote path mapping %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update remote path mapping %d: %w", id, err)
	}

	return ptrResult[RemotePathMapping](resp)
}

// Delete removes a remote path mapping by ID.
func (s *RemotePathMappingService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/remotepathmapping/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete remote path mapping %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete remote path mapping %d: %w", id, err)
	}

	return nil
}
