package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// ExclusionsService provides methods for the /exclusions endpoint.
type ExclusionsService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ImportListExclusion is a movie excluded from import list monitoring.
type ImportListExclusion struct {
	ID         int     `json:"id"`
	TmdbID     int     `json:"tmdbId"`
	MovieTitle *string `json:"movieTitle,omitempty"`
	MovieYear  int     `json:"movieYear"`
}

// ExclusionsBulkRequest carries IDs for a bulk delete operation.
type ExclusionsBulkRequest struct {
	IDs []int `json:"ids,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all import list exclusions.
func (s *ExclusionsService) List(ctx context.Context) ([]ImportListExclusion, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ImportListExclusion{}).
		Get("/api/v3/exclusions")
	if err != nil {
		return nil, fmt.Errorf("radarr: list exclusions: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list exclusions: %w", err)
	}

	return derefResult[[]ImportListExclusion](resp)
}

// Get returns a single exclusion by ID.
func (s *ExclusionsService) Get(ctx context.Context, id int) (*ImportListExclusion, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportListExclusion{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/exclusions/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get exclusion %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get exclusion %d: %w", id, err)
	}

	return ptrResult[ImportListExclusion](resp)
}

// ---------------------------------------------------------------------------
// ListPaged
// ---------------------------------------------------------------------------

// ListExclusionsPagedOption is a functional option for ExclusionsService.ListPaged.
type ListExclusionsPagedOption func(*resty.Request)

// WithExclusionsPage requests a specific page number.
func WithExclusionsPage(page int) ListExclusionsPagedOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithExclusionsPageSize sets the number of records per page.
func WithExclusionsPageSize(n int) ListExclusionsPagedOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithExclusionsSortKey sets the field used to sort results.
func WithExclusionsSortKey(key string) ListExclusionsPagedOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithExclusionsSortDirection sets the sort order.
func WithExclusionsSortDirection(dir SortDirection) ListExclusionsPagedOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// ListPaged returns a single page of exclusions.
func (s *ExclusionsService) ListPaged(ctx context.Context, opts ...ListExclusionsPagedOption) (*PagedResult[ImportListExclusion], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[ImportListExclusion]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/exclusions/paged")
	if err != nil {
		return nil, fmt.Errorf("radarr: list exclusions paged: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list exclusions paged: %w", err)
	}

	return ptrResult[PagedResult[ImportListExclusion]](resp)
}

// ListPagedWithPagination iterates over all exclusion pages automatically.
func (s *ExclusionsService) ListPagedWithPagination(ctx context.Context, opts ...ListExclusionsPagedOption) iter.Seq2[PagedResult[ImportListExclusion], error] {
	return paginate(ctx, s.ListPaged, opts...)
}

// ---------------------------------------------------------------------------
// Create / CreateBulk / Delete / DeleteBulk
// ---------------------------------------------------------------------------

// Create adds a new import list exclusion.
func (s *ExclusionsService) Create(ctx context.Context, body ImportListExclusion) (*ImportListExclusion, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportListExclusion{}).
		SetBody(body).
		Post("/api/v3/exclusions")
	if err != nil {
		return nil, fmt.Errorf("radarr: create exclusion: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create exclusion: %w", err)
	}

	return ptrResult[ImportListExclusion](resp)
}

// CreateBulk adds multiple import list exclusions in a single request.
func (s *ExclusionsService) CreateBulk(ctx context.Context, body []ImportListExclusion) ([]ImportListExclusion, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ImportListExclusion{}).
		SetBody(body).
		Post("/api/v3/exclusions/bulk")
	if err != nil {
		return nil, fmt.Errorf("radarr: create exclusions bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create exclusions bulk: %w", err)
	}

	return derefResult[[]ImportListExclusion](resp)
}

// Delete removes a single import list exclusion by ID.
func (s *ExclusionsService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/exclusions/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete exclusion %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete exclusion %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple import list exclusions in a single request.
func (s *ExclusionsService) DeleteBulk(ctx context.Context, body ExclusionsBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/exclusions/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete exclusions bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete exclusions bulk: %w", err)
	}

	return nil
}
