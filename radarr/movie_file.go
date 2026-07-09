package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// MovieFileService provides methods for the /moviefile endpoint.
type MovieFileService service

// MovieFileBulkRequest is the body for bulk delete and editor update operations.
type MovieFileBulkRequest struct {
	MovieFileIDs []int         `json:"movieFileIds,omitempty"`
	Languages    []Language    `json:"languages,omitempty"`
	Quality      *QualityModel `json:"quality,omitempty"`
	Edition      *string       `json:"edition,omitempty"`
	ReleaseGroup *string       `json:"releaseGroup,omitempty"`
	SceneName    *string       `json:"sceneName,omitempty"`
	IndexerFlags *int          `json:"indexerFlags,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListMovieFilesOption is a functional option for MovieFileService.List.
type ListMovieFilesOption func(*resty.Request)

// WithMovieFileMovieIDs filters results to files belonging to the given movie IDs.
func WithMovieFileMovieIDs(ids ...int) ListMovieFilesOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("movieId", strconv.Itoa(id))
		}
	}
}

// WithMovieFileIDs filters results to the given movie file IDs.
func WithMovieFileIDs(ids ...int) ListMovieFilesOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("movieFileIds", strconv.Itoa(id))
		}
	}
}

// List returns movie files matching the provided filters.
func (s *MovieFileService) List(ctx context.Context, opts ...ListMovieFilesOption) ([]MovieFile, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]MovieFile{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/moviefile")
	if err != nil {
		return nil, fmt.Errorf("radarr: list movie files: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list movie files: %w", err)
	}

	return derefResult[[]MovieFile](resp)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

// Get returns a single movie file by ID.
func (s *MovieFileService) Get(ctx context.Context, id int) (*MovieFile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MovieFile{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/moviefile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get movie file %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get movie file %d: %w", id, err)
	}

	return ptrResult[MovieFile](resp)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// Update replaces a movie file's metadata by ID.
func (s *MovieFileService) Update(ctx context.Context, id int, body MovieFile) (*MovieFile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MovieFile{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/moviefile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update movie file %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update movie file %d: %w", id, err)
	}

	return ptrResult[MovieFile](resp)
}

// ---------------------------------------------------------------------------
// UpdateBulk
// ---------------------------------------------------------------------------

// UpdateBulk replaces multiple movie files in a single request.
func (s *MovieFileService) UpdateBulk(ctx context.Context, files []MovieFile) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(files).
		Put("/api/v3/moviefile/bulk")
	if err != nil {
		return fmt.Errorf("radarr: update movie files bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: update movie files bulk: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// UpdateEditor
// ---------------------------------------------------------------------------

// UpdateEditor applies common fields to a set of movie files identified by IDs.
func (s *MovieFileService) UpdateEditor(ctx context.Context, body MovieFileBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Put("/api/v3/moviefile/editor")
	if err != nil {
		return fmt.Errorf("radarr: update movie file editor: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: update movie file editor: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// Delete removes a movie file by ID.
func (s *MovieFileService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/moviefile/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete movie file %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete movie file %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// DeleteBulk
// ---------------------------------------------------------------------------

// DeleteBulk removes multiple movie files in a single request.
func (s *MovieFileService) DeleteBulk(ctx context.Context, body MovieFileBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/moviefile/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete movie files bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete movie files bulk: %w", err)
	}

	return nil
}
