package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// AltTitleService provides methods for the /alttitle endpoint.
type AltTitleService service

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListAltTitlesOption is a functional option for AltTitleService.List.
type ListAltTitlesOption func(*resty.Request)

// WithAltTitleMovieID filters alternative titles to a specific movie ID.
func WithAltTitleMovieID(id int) ListAltTitlesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// WithAltTitleMovieMetadataID filters alternative titles to a specific movie metadata ID.
func WithAltTitleMovieMetadataID(id int) ListAltTitlesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieMetadataId", strconv.Itoa(id))
	}
}

// List returns alternative titles, optionally filtered by movie or metadata ID.
func (s *AltTitleService) List(ctx context.Context, opts ...ListAltTitlesOption) ([]AlternativeTitle, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]AlternativeTitle{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/alttitle")
	if err != nil {
		return nil, fmt.Errorf("radarr: list alt titles: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list alt titles: %w", err)
	}

	return derefResult[[]AlternativeTitle](resp)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

// Get returns a single alternative title by ID.
func (s *AltTitleService) Get(ctx context.Context, id int) (*AlternativeTitle, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&AlternativeTitle{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/alttitle/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get alt title %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get alt title %d: %w", id, err)
	}

	return ptrResult[AlternativeTitle](resp)
}
