package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// RenameService provides methods for the /rename endpoint.
type RenameService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// RenameMovie describes the old and new paths for a file that would be renamed.
type RenameMovie struct {
	ID           int     `json:"id"`
	MovieID      int     `json:"movieId"`
	MovieFileID  int     `json:"movieFileId"`
	ExistingPath *string `json:"existingPath,omitempty"`
	NewPath      *string `json:"newPath,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListRenameOption is a functional option for RenameService.List.
type ListRenameOption func(*resty.Request)

// WithRenameMovieID filters rename candidates to a specific movie ID.
func WithRenameMovieID(id int) ListRenameOption {
	return func(r *resty.Request) {
		r.QueryParam.Add("movieId", strconv.Itoa(id))
	}
}

// List returns files that would be renamed according to the current naming configuration.
func (s *RenameService) List(ctx context.Context, opts ...ListRenameOption) ([]RenameMovie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]RenameMovie{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/rename")
	if err != nil {
		return nil, fmt.Errorf("radarr: list rename: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list rename: %w", err)
	}

	return derefResult[[]RenameMovie](resp)
}
