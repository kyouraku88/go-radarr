package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// HistoryService provides methods for the /history endpoint.
type HistoryService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// MovieHistoryEventType classifies a history entry.
type MovieHistoryEventType string

// Movie history event type values.
const (
	MovieHistoryEventTypeUnknown                MovieHistoryEventType = "unknown"
	MovieHistoryEventTypeGrabbed                MovieHistoryEventType = "grabbed"
	MovieHistoryEventTypeDownloadFolderImported MovieHistoryEventType = "downloadFolderImported"
	MovieHistoryEventTypeDownloadFailed         MovieHistoryEventType = "downloadFailed"
	MovieHistoryEventTypeMovieFileDeleted       MovieHistoryEventType = "movieFileDeleted"
	MovieHistoryEventTypeMovieFolderImported    MovieHistoryEventType = "movieFolderImported"
	MovieHistoryEventTypeMovieFileRenamed       MovieHistoryEventType = "movieFileRenamed"
	MovieHistoryEventTypeDownloadIgnored        MovieHistoryEventType = "downloadIgnored"
)

// HistoryRecord is a single entry in the download history.
type HistoryRecord struct {
	ID                  int                   `json:"id"`
	MovieID             int                   `json:"movieId"`
	SourceTitle         *string               `json:"sourceTitle,omitempty"`
	Languages           []Language            `json:"languages,omitempty"`
	Quality             QualityModel          `json:"quality"`
	CustomFormats       []CustomFormat        `json:"customFormats,omitempty"`
	CustomFormatScore   int                   `json:"customFormatScore"`
	QualityCutoffNotMet bool                  `json:"qualityCutoffNotMet"`
	Date                time.Time             `json:"date"`
	DownloadID          *string               `json:"downloadId,omitempty"`
	EventType           MovieHistoryEventType `json:"eventType,omitempty"`
	Data                map[string]*string    `json:"data,omitempty"`
	Movie               *Movie                `json:"movie,omitempty"`
}

// ---------------------------------------------------------------------------
// List (paginated)
// ---------------------------------------------------------------------------

// ListHistoryOption is a functional option for HistoryService.List.
type ListHistoryOption func(*resty.Request)

// WithHistoryPage requests a specific page number.
func WithHistoryPage(n int) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(n))
	}
}

// WithHistoryPageSize sets the number of records per page.
func WithHistoryPageSize(n int) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithHistorySortKey sets the field used to sort results.
func WithHistorySortKey(key string) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithHistorySortDirection sets the sort order.
func WithHistorySortDirection(dir SortDirection) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithHistoryIncludeMovie embeds movie details in each record when v is true.
func WithHistoryIncludeMovie(v bool) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeMovie", strconv.FormatBool(v))
	}
}

// WithHistoryEventTypeIDs filters by numeric event type IDs (as accepted by the paginated endpoint).
func WithHistoryEventTypeIDs(ids ...int) ListHistoryOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("eventType", strconv.Itoa(id))
		}
	}
}

// WithHistoryDownloadID filters records by download client ID.
func WithHistoryDownloadID(id string) ListHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("downloadId", id)
	}
}

// WithHistoryMovieIDs filters records to specific movie IDs.
func WithHistoryMovieIDs(ids ...int) ListHistoryOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("movieIds", strconv.Itoa(id))
		}
	}
}

// WithHistoryLanguageIDs filters records by language IDs.
func WithHistoryLanguageIDs(ids ...int) ListHistoryOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("languages", strconv.Itoa(id))
		}
	}
}

// WithHistoryQualityIDs filters records by quality IDs.
func WithHistoryQualityIDs(ids ...int) ListHistoryOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("quality", strconv.Itoa(id))
		}
	}
}

// List returns a single page of history records.
func (s *HistoryService) List(ctx context.Context, opts ...ListHistoryOption) (*PagedResult[HistoryRecord], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[HistoryRecord]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/history")
	if err != nil {
		return nil, fmt.Errorf("radarr: list history: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list history: %w", err)
	}

	return ptrResult[PagedResult[HistoryRecord]](resp)
}

// ListWithPagination iterates over all history pages automatically.
// Use range to consume pages:
//
//	for page, err := range client.History.ListWithPagination(ctx, WithHistoryPageSize(100)) {
//	    if err != nil { ... }
//	    for _, record := range page.Records { ... }
//	}
func (s *HistoryService) ListWithPagination(ctx context.Context, opts ...ListHistoryOption) iter.Seq2[PagedResult[HistoryRecord], error] {
	return paginate(ctx, s.List, opts...)
}

// ---------------------------------------------------------------------------
// GetByMovie
// ---------------------------------------------------------------------------

// GetMovieHistoryOption is a functional option for HistoryService.GetByMovie.
type GetMovieHistoryOption func(*resty.Request)

// WithMovieHistoryEventType filters history to a specific event type.
func WithMovieHistoryEventType(t MovieHistoryEventType) GetMovieHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("eventType", string(t))
	}
}

// WithMovieHistoryIncludeMovie embeds movie details in each record when v is true.
func WithMovieHistoryIncludeMovie(v bool) GetMovieHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeMovie", strconv.FormatBool(v))
	}
}

// GetByMovie returns the download history for a single movie.
func (s *HistoryService) GetByMovie(ctx context.Context, movieID int, opts ...GetMovieHistoryOption) ([]HistoryRecord, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]HistoryRecord{}).
		SetQueryParam("movieId", strconv.Itoa(movieID))
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/history/movie")
	if err != nil {
		return nil, fmt.Errorf("radarr: get history for movie %d: %w", movieID, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get history for movie %d: %w", movieID, err)
	}

	return derefResult[[]HistoryRecord](resp)
}

// ---------------------------------------------------------------------------
// Since
// ---------------------------------------------------------------------------

// ListSinceHistoryOption is a functional option for HistoryService.Since.
type ListSinceHistoryOption func(*resty.Request)

// WithSinceEventType filters history to a specific event type.
func WithSinceEventType(t MovieHistoryEventType) ListSinceHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("eventType", string(t))
	}
}

// WithSinceIncludeMovie embeds movie details in each record when v is true.
func WithSinceIncludeMovie(v bool) ListSinceHistoryOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeMovie", strconv.FormatBool(v))
	}
}

// Since returns all history records after date.
func (s *HistoryService) Since(ctx context.Context, date time.Time, opts ...ListSinceHistoryOption) ([]HistoryRecord, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]HistoryRecord{}).
		SetQueryParam("date", date.UTC().Format(time.RFC3339))
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/history/since")
	if err != nil {
		return nil, fmt.Errorf("radarr: list history since %s: %w", date.Format(time.DateOnly), err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list history since %s: %w", date.Format(time.DateOnly), err)
	}

	return derefResult[[]HistoryRecord](resp)
}

// ---------------------------------------------------------------------------
// MarkFailed
// ---------------------------------------------------------------------------

// MarkFailed marks a history record as failed, triggering a re-search.
func (s *HistoryService) MarkFailed(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Post("/api/v3/history/failed/{id}")
	if err != nil {
		return fmt.Errorf("radarr: mark history %d failed: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: mark history %d failed: %w", id, err)
	}

	return nil
}
