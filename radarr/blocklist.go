package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// BlocklistService provides methods for the /blocklist endpoint.
type BlocklistService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// BlocklistRecord is a release that has been added to the blocklist.
type BlocklistRecord struct {
	ID            int              `json:"id"`
	MovieID       int              `json:"movieId"`
	SourceTitle   *string          `json:"sourceTitle,omitempty"`
	Languages     []Language       `json:"languages,omitempty"`
	Quality       QualityModel     `json:"quality"`
	CustomFormats []CustomFormat   `json:"customFormats,omitempty"`
	Date          time.Time        `json:"date"`
	Protocol      DownloadProtocol `json:"protocol,omitempty"`
	Indexer       *string          `json:"indexer,omitempty"`
	Message       *string          `json:"message,omitempty"`
	Movie         *Movie           `json:"movie,omitempty"`
}

// BlocklistBulkRequest carries a list of blocklist IDs for bulk deletion.
type BlocklistBulkRequest struct {
	IDs []int `json:"ids,omitempty"`
}

// ---------------------------------------------------------------------------
// List (paged)
// ---------------------------------------------------------------------------

// ListBlocklistOption is a functional option for BlocklistService.List.
type ListBlocklistOption func(*resty.Request)

// WithBlocklistPage requests a specific page number.
func WithBlocklistPage(page int) ListBlocklistOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithBlocklistPageSize sets the number of records per page.
func WithBlocklistPageSize(n int) ListBlocklistOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithBlocklistSortKey sets the field used to sort results.
func WithBlocklistSortKey(key string) ListBlocklistOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithBlocklistSortDirection sets the sort order.
func WithBlocklistSortDirection(dir SortDirection) ListBlocklistOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithBlocklistMovieIDs filters results to specific movie IDs.
func WithBlocklistMovieIDs(ids ...int) ListBlocklistOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("movieIds", strconv.Itoa(id))
		}
	}
}

// WithBlocklistProtocols filters results to specific protocols.
func WithBlocklistProtocols(protocols ...DownloadProtocol) ListBlocklistOption {
	return func(r *resty.Request) {
		for _, p := range protocols {
			r.QueryParam.Add("protocols", string(p))
		}
	}
}

// List returns a single page of blocklist records.
func (s *BlocklistService) List(ctx context.Context, opts ...ListBlocklistOption) (*PagedResult[BlocklistRecord], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[BlocklistRecord]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/blocklist")
	if err != nil {
		return nil, fmt.Errorf("radarr: list blocklist: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list blocklist: %w", err)
	}

	return ptrResult[PagedResult[BlocklistRecord]](resp)
}

// ListWithPagination iterates over all blocklist pages automatically.
func (s *BlocklistService) ListWithPagination(ctx context.Context, opts ...ListBlocklistOption) iter.Seq2[PagedResult[BlocklistRecord], error] {
	return paginate(ctx, s.List, opts...)
}

// ---------------------------------------------------------------------------
// ListByMovie
// ---------------------------------------------------------------------------

// ListByMovie returns blocklist records for a specific movie.
func (s *BlocklistService) ListByMovie(ctx context.Context, movieID int) ([]BlocklistRecord, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]BlocklistRecord{}).
		SetQueryParam("movieId", strconv.Itoa(movieID)).
		Get("/api/v3/blocklist/movie")
	if err != nil {
		return nil, fmt.Errorf("radarr: list blocklist by movie %d: %w", movieID, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list blocklist by movie %d: %w", movieID, err)
	}

	return derefResult[[]BlocklistRecord](resp)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// Delete removes a single blocklist entry by ID.
func (s *BlocklistService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/blocklist/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete blocklist entry %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete blocklist entry %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple blocklist entries in a single request.
func (s *BlocklistService) DeleteBulk(ctx context.Context, body BlocklistBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/blocklist/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete blocklist bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete blocklist bulk: %w", err)
	}

	return nil
}
