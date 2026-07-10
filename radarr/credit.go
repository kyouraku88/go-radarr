package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// CreditService provides methods for the /credit endpoint.
type CreditService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Credit represents a cast or crew credit for a movie.
type Credit struct {
	ID              int          `json:"id"`
	PersonName      *string      `json:"personName,omitempty"`
	CreditTmdbID    *string      `json:"creditTmdbId,omitempty"`
	PersonTmdbID    int          `json:"personTmdbId"`
	MovieMetadataID int          `json:"movieMetadataId"`
	Images          []MediaCover `json:"images,omitempty"`
	Department      *string      `json:"department,omitempty"`
	Job             *string      `json:"job,omitempty"`
	Character       *string      `json:"character,omitempty"`
	Order           int          `json:"order"`
	Type            CreditType   `json:"type,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListCreditsOption is a functional option for CreditService.List.
type ListCreditsOption func(*resty.Request)

// WithCreditMovieID filters credits to a specific movie ID.
func WithCreditMovieID(id int) ListCreditsOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// WithCreditMovieMetadataID filters credits to a specific movie metadata ID.
func WithCreditMovieMetadataID(id int) ListCreditsOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieMetadataId", strconv.Itoa(id))
	}
}

// List returns movie credits, optionally filtered by movie or metadata ID.
func (s *CreditService) List(ctx context.Context, opts ...ListCreditsOption) ([]Credit, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Credit{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/credit")
	if err != nil {
		return nil, fmt.Errorf("radarr: list credits: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list credits: %w", err)
	}

	return derefResult[[]Credit](resp)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

// Get returns a single credit by ID.
func (s *CreditService) Get(ctx context.Context, id int) (*Credit, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Credit{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/credit/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get credit %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get credit %d: %w", id, err)
	}

	return ptrResult[Credit](resp)
}
