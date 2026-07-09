package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// TagService provides methods for the /tag endpoint.
type TagService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Tag is a simple label that can be attached to movies, profiles, and other resources.
type Tag struct {
	ID    int     `json:"id"`
	Label *string `json:"label,omitempty"`
}

// TagDetails extends Tag with lists of all resources that reference it.
type TagDetails struct {
	ID                int     `json:"id"`
	Label             *string `json:"label,omitempty"`
	DelayProfileIDs   []int   `json:"delayProfileIds,omitempty"`
	ImportListIDs     []int   `json:"importListIds,omitempty"`
	NotificationIDs   []int   `json:"notificationIds,omitempty"`
	ReleaseProfileIDs []int   `json:"releaseProfileIds,omitempty"`
	IndexerIDs        []int   `json:"indexerIds,omitempty"`
	DownloadClientIDs []int   `json:"downloadClientIds,omitempty"`
	AutoTagIDs        []int   `json:"autoTagIds,omitempty"`
	MovieIDs          []int   `json:"movieIds,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all tags.
func (s *TagService) List(ctx context.Context) ([]Tag, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Tag{}).
		Get("/api/v3/tag")
	if err != nil {
		return nil, fmt.Errorf("radarr: list tags: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list tags: %w", err)
	}

	return derefResult[[]Tag](resp)
}

// Get returns a single tag by ID.
func (s *TagService) Get(ctx context.Context, id int) (*Tag, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Tag{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/tag/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get tag %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get tag %d: %w", id, err)
	}

	return ptrResult[Tag](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new tag.
func (s *TagService) Create(ctx context.Context, body Tag) (*Tag, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Tag{}).
		SetBody(body).
		Post("/api/v3/tag")
	if err != nil {
		return nil, fmt.Errorf("radarr: create tag: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create tag: %w", err)
	}

	return ptrResult[Tag](resp)
}

// Update replaces a tag by ID.
func (s *TagService) Update(ctx context.Context, id int, body Tag) (*Tag, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Tag{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/tag/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update tag %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update tag %d: %w", id, err)
	}

	return ptrResult[Tag](resp)
}

// Delete removes a tag by ID.
func (s *TagService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/tag/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete tag %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete tag %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Details
// ---------------------------------------------------------------------------

// ListDetails returns all tags with their associated resource IDs.
func (s *TagService) ListDetails(ctx context.Context) ([]TagDetails, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]TagDetails{}).
		Get("/api/v3/tag/detail")
	if err != nil {
		return nil, fmt.Errorf("radarr: list tag details: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list tag details: %w", err)
	}

	return derefResult[[]TagDetails](resp)
}

// GetDetail returns a single tag with its associated resource IDs by tag ID.
func (s *TagService) GetDetail(ctx context.Context, id int) (*TagDetails, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&TagDetails{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/tag/detail/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get tag detail %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get tag detail %d: %w", id, err)
	}

	return ptrResult[TagDetails](resp)
}
