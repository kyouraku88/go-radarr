package radarr

import (
	"context"
	"fmt"
)

// MovieEditorService provides bulk edit and delete methods for the /movie/editor endpoint.
type MovieEditorService service

// ApplyTags controls how tags are merged with existing ones during a bulk edit.
type ApplyTags string

// Apply tags values.
const (
	ApplyTagsAdd     ApplyTags = "add"
	ApplyTagsRemove  ApplyTags = "remove"
	ApplyTagsReplace ApplyTags = "replace"
)

// MovieEditorRequest is the body for bulk-edit and bulk-delete operations.
type MovieEditorRequest struct {
	MovieIDs            []int           `json:"movieIds,omitempty"`
	Monitored           *bool           `json:"monitored,omitempty"`
	QualityProfileID    *int            `json:"qualityProfileId,omitempty"`
	MinimumAvailability MovieStatusType `json:"minimumAvailability,omitempty"`
	RootFolderPath      *string         `json:"rootFolderPath,omitempty"`
	Tags                []int           `json:"tags,omitempty"`
	ApplyTags           ApplyTags       `json:"applyTags,omitempty"`
	MoveFiles           bool            `json:"moveFiles,omitempty"`
	DeleteFiles         bool            `json:"deleteFiles,omitempty"`
	AddImportExclusion  bool            `json:"addImportExclusion,omitempty"`
}

// Edit performs a bulk update of movies matching the request.
func (s *MovieEditorService) Edit(ctx context.Context, body MovieEditorRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Put("/api/v3/movie/editor")
	if err != nil {
		return fmt.Errorf("radarr: edit movies: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: edit movies: %w", err)
	}

	return nil
}

// Delete performs a bulk deletion of movies matching the request.
func (s *MovieEditorService) Delete(ctx context.Context, body MovieEditorRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/movie/editor")
	if err != nil {
		return fmt.Errorf("radarr: delete movies: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete movies: %w", err)
	}

	return nil
}
