package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// RootFolderService provides methods for the /rootfolder endpoint.
type RootFolderService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// UnmappedFolder is a subfolder inside a root folder that has no matched movie.
type UnmappedFolder struct {
	Name         *string `json:"name,omitempty"`
	Path         *string `json:"path,omitempty"`
	RelativePath *string `json:"relativePath,omitempty"`
}

// RootFolder is a top-level media directory configured in Radarr.
type RootFolder struct {
	ID              int              `json:"id"`
	Path            *string          `json:"path,omitempty"`
	Accessible      bool             `json:"accessible"`
	FreeSpace       *int64           `json:"freeSpace,omitempty"`
	UnmappedFolders []UnmappedFolder `json:"unmappedFolders,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all configured root folders.
func (s *RootFolderService) List(ctx context.Context) ([]RootFolder, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]RootFolder{}).
		Get("/api/v3/rootfolder")
	if err != nil {
		return nil, fmt.Errorf("radarr: list root folders: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list root folders: %w", err)
	}

	return derefResult[[]RootFolder](resp)
}

// Get returns a single root folder by ID.
func (s *RootFolderService) Get(ctx context.Context, id int) (*RootFolder, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&RootFolder{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/rootfolder/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get root folder %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get root folder %d: %w", id, err)
	}

	return ptrResult[RootFolder](resp)
}

// ---------------------------------------------------------------------------
// Create / Delete
// ---------------------------------------------------------------------------

// Create adds a new root folder.
func (s *RootFolderService) Create(ctx context.Context, body RootFolder) (*RootFolder, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&RootFolder{}).
		SetBody(body).
		Post("/api/v3/rootfolder")
	if err != nil {
		return nil, fmt.Errorf("radarr: create root folder: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create root folder: %w", err)
	}

	return ptrResult[RootFolder](resp)
}

// Delete removes a root folder by ID.
func (s *RootFolderService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/rootfolder/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete root folder %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete root folder %d: %w", id, err)
	}

	return nil
}
