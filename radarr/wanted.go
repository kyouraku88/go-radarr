package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// WantedService provides methods for the /wanted endpoint.
type WantedService service

// ---------------------------------------------------------------------------
// Cutoff
// ---------------------------------------------------------------------------

// ListWantedCutoffOption is a functional option for WantedService.ListCutoff.
type ListWantedCutoffOption func(*resty.Request)

// WithWantedCutoffPage requests a specific page number.
func WithWantedCutoffPage(page int) ListWantedCutoffOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithWantedCutoffPageSize sets the number of records per page.
func WithWantedCutoffPageSize(n int) ListWantedCutoffOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithWantedCutoffSortKey sets the field used to sort results.
func WithWantedCutoffSortKey(key string) ListWantedCutoffOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithWantedCutoffSortDirection sets the sort order.
func WithWantedCutoffSortDirection(dir SortDirection) ListWantedCutoffOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithWantedCutoffMonitored filters to monitored movies when v is true.
func WithWantedCutoffMonitored(v bool) ListWantedCutoffOption {
	return func(r *resty.Request) {
		r.SetQueryParam("monitored", strconv.FormatBool(v))
	}
}

// ListCutoff returns a single page of movies that have not met their quality cutoff.
func (s *WantedService) ListCutoff(ctx context.Context, opts ...ListWantedCutoffOption) (*PagedResult[Movie], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[Movie]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/wanted/cutoff")
	if err != nil {
		return nil, fmt.Errorf("radarr: list wanted cutoff: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list wanted cutoff: %w", err)
	}

	return ptrResult[PagedResult[Movie]](resp)
}

// ListCutoffWithPagination iterates over all cutoff pages automatically.
func (s *WantedService) ListCutoffWithPagination(ctx context.Context, opts ...ListWantedCutoffOption) iter.Seq2[PagedResult[Movie], error] {
	return paginate(ctx, s.ListCutoff, opts...)
}

// ---------------------------------------------------------------------------
// Missing
// ---------------------------------------------------------------------------

// ListWantedMissingOption is a functional option for WantedService.ListMissing.
type ListWantedMissingOption func(*resty.Request)

// WithWantedMissingPage requests a specific page number.
func WithWantedMissingPage(page int) ListWantedMissingOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithWantedMissingPageSize sets the number of records per page.
func WithWantedMissingPageSize(n int) ListWantedMissingOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithWantedMissingSortKey sets the field used to sort results.
func WithWantedMissingSortKey(key string) ListWantedMissingOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithWantedMissingSortDirection sets the sort order.
func WithWantedMissingSortDirection(dir SortDirection) ListWantedMissingOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithWantedMissingMonitored filters to monitored movies when v is true.
func WithWantedMissingMonitored(v bool) ListWantedMissingOption {
	return func(r *resty.Request) {
		r.SetQueryParam("monitored", strconv.FormatBool(v))
	}
}

// ListMissing returns a single page of monitored movies without a downloaded file.
func (s *WantedService) ListMissing(ctx context.Context, opts ...ListWantedMissingOption) (*PagedResult[Movie], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[Movie]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/wanted/missing")
	if err != nil {
		return nil, fmt.Errorf("radarr: list wanted missing: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list wanted missing: %w", err)
	}

	return ptrResult[PagedResult[Movie]](resp)
}

// ListMissingWithPagination iterates over all missing pages automatically.
func (s *WantedService) ListMissingWithPagination(ctx context.Context, opts ...ListWantedMissingOption) iter.Seq2[PagedResult[Movie], error] {
	return paginate(ctx, s.ListMissing, opts...)
}
