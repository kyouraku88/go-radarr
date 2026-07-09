package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// QueueService provides methods for the /queue endpoint.
type QueueService service

// ---------------------------------------------------------------------------
// Queue-specific types
// ---------------------------------------------------------------------------

// QueueStatus represents the download status of a queue item.
type QueueStatus string

// Queue status values.
const (
	QueueStatusUnknown                   QueueStatus = "unknown"
	QueueStatusQueued                    QueueStatus = "queued"
	QueueStatusPaused                    QueueStatus = "paused"
	QueueStatusDownloading               QueueStatus = "downloading"
	QueueStatusCompleted                 QueueStatus = "completed"
	QueueStatusFailed                    QueueStatus = "failed"
	QueueStatusWarning                   QueueStatus = "warning"
	QueueStatusDelay                     QueueStatus = "delay"
	QueueStatusDownloadClientUnavailable QueueStatus = "downloadClientUnavailable"
	QueueStatusFallback                  QueueStatus = "fallback"
)

// TrackedDownloadStatus indicates the health of a tracked download.
type TrackedDownloadStatus string

// Tracked download status values.
const (
	TrackedDownloadStatusOk      TrackedDownloadStatus = "ok"
	TrackedDownloadStatusWarning TrackedDownloadStatus = "warning"
	TrackedDownloadStatusError   TrackedDownloadStatus = "error"
)

// TrackedDownloadState describes the lifecycle state of a tracked download.
type TrackedDownloadState string

// Tracked download state values.
const (
	TrackedDownloadStateDownloading   TrackedDownloadState = "downloading"
	TrackedDownloadStateImportBlocked TrackedDownloadState = "importBlocked"
	TrackedDownloadStateImportPending TrackedDownloadState = "importPending"
	TrackedDownloadStateImporting     TrackedDownloadState = "importing"
	TrackedDownloadStateImported      TrackedDownloadState = "imported"
	TrackedDownloadStateFailedPending TrackedDownloadState = "failedPending"
	TrackedDownloadStateFailed        TrackedDownloadState = "failed"
	TrackedDownloadStateIgnored       TrackedDownloadState = "ignored"
)

// TrackedDownloadStatusMessage is a diagnostic message attached to a queue record.
type TrackedDownloadStatusMessage struct {
	Title    *string  `json:"title,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

// QueueRecord represents a single item in the download queue.
type QueueRecord struct {
	ID                                  int                            `json:"id"`
	MovieID                             *int                           `json:"movieId,omitempty"`
	Movie                               *Movie                         `json:"movie,omitempty"`
	Languages                           []Language                     `json:"languages,omitempty"`
	Quality                             QualityModel                   `json:"quality"`
	CustomFormats                       []CustomFormat                 `json:"customFormats,omitempty"`
	CustomFormatScore                   int                            `json:"customFormatScore"`
	Size                                float64                        `json:"size"`
	Title                               *string                        `json:"title,omitempty"`
	EstimatedCompletionTime             *time.Time                     `json:"estimatedCompletionTime,omitempty"`
	Added                               *time.Time                     `json:"added,omitempty"`
	Status                              QueueStatus                    `json:"status,omitempty"`
	TrackedDownloadStatus               TrackedDownloadStatus          `json:"trackedDownloadStatus,omitempty"`
	TrackedDownloadState                TrackedDownloadState           `json:"trackedDownloadState,omitempty"`
	StatusMessages                      []TrackedDownloadStatusMessage `json:"statusMessages,omitempty"`
	ErrorMessage                        *string                        `json:"errorMessage,omitempty"`
	DownloadID                          *string                        `json:"downloadId,omitempty"`
	Protocol                            DownloadProtocol               `json:"protocol,omitempty"`
	DownloadClient                      *string                        `json:"downloadClient,omitempty"`
	DownloadClientHasPostImportCategory bool                           `json:"downloadClientHasPostImportCategory"`
	Indexer                             *string                        `json:"indexer,omitempty"`
	OutputPath                          *string                        `json:"outputPath,omitempty"`
}

// QueueBulkRequest carries a list of queue item IDs for bulk operations.
type QueueBulkRequest struct {
	IDs []int `json:"ids,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListQueueOption is a functional option for QueueService.List.
type ListQueueOption func(*resty.Request)

// WithQueuePage requests a specific page number.
func WithQueuePage(page int) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithQueuePageSize sets the number of records per page.
func WithQueuePageSize(n int) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithQueueSortKey sets the field used to sort results.
func WithQueueSortKey(key string) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithQueueSortDirection sets the sort order.
func WithQueueSortDirection(dir SortDirection) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithQueueIncludeUnknownMovieItems includes items without a matched movie when v is true.
func WithQueueIncludeUnknownMovieItems(v bool) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeUnknownMovieItems", strconv.FormatBool(v))
	}
}

// WithQueueIncludeMovie embeds movie details in each record when v is true.
func WithQueueIncludeMovie(v bool) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeMovie", strconv.FormatBool(v))
	}
}

// WithQueueMovieIDs filters results to specific movie IDs.
func WithQueueMovieIDs(ids ...int) ListQueueOption {
	return func(r *resty.Request) {
		for _, id := range ids {
			r.QueryParam.Add("movieIds", strconv.Itoa(id))
		}
	}
}

// WithQueueProtocol filters results by download protocol.
func WithQueueProtocol(p DownloadProtocol) ListQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("protocol", string(p))
	}
}

// WithQueueStatuses filters results to items with the given statuses.
func WithQueueStatuses(statuses ...QueueStatus) ListQueueOption {
	return func(r *resty.Request) {
		for _, s := range statuses {
			r.QueryParam.Add("status", string(s))
		}
	}
}

// List returns a single page of queue records.
func (s *QueueService) List(ctx context.Context, opts ...ListQueueOption) (*PagedResult[QueueRecord], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[QueueRecord]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/queue")
	if err != nil {
		return nil, fmt.Errorf("radarr: list queue: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list queue: %w", err)
	}

	return ptrResult[PagedResult[QueueRecord]](resp)
}

// ListWithPagination iterates over all pages automatically.
// Use range to consume pages:
//
//	for page, err := range client.Queue.ListWithPagination(ctx, WithQueuePageSize(100)) {
//	    if err != nil { ... }
//	    for _, record := range page.Records { ... }
//	}
func (s *QueueService) ListWithPagination(ctx context.Context, opts ...ListQueueOption) iter.Seq2[PagedResult[QueueRecord], error] {
	return paginate(ctx, s.List, opts...)
}

// ---------------------------------------------------------------------------
// Details
// ---------------------------------------------------------------------------

// ListQueueDetailsOption is a functional option for QueueService.ListDetails.
type ListQueueDetailsOption func(*resty.Request)

// WithQueueDetailsMovieID filters detail records to a specific movie ID.
func WithQueueDetailsMovieID(id int) ListQueueDetailsOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// WithQueueDetailsIncludeMovie embeds movie details in each record when v is true.
func WithQueueDetailsIncludeMovie(v bool) ListQueueDetailsOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeMovie", strconv.FormatBool(v))
	}
}

// ListDetails returns queue records from the /queue/details endpoint.
func (s *QueueService) ListDetails(ctx context.Context, opts ...ListQueueDetailsOption) ([]QueueRecord, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]QueueRecord{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/queue/details")
	if err != nil {
		return nil, fmt.Errorf("radarr: list queue details: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list queue details: %w", err)
	}

	return derefResult[[]QueueRecord](resp)
}

// ---------------------------------------------------------------------------
// Status
// ---------------------------------------------------------------------------

// QueueStatusSummary is a summary of the current queue counters.
type QueueStatusSummary struct {
	ID              int  `json:"id"`
	TotalCount      int  `json:"totalCount"`
	Count           int  `json:"count"`
	UnknownCount    int  `json:"unknownCount"`
	Errors          bool `json:"errors"`
	Warnings        bool `json:"warnings"`
	UnknownErrors   bool `json:"unknownErrors"`
	UnknownWarnings bool `json:"unknownWarnings"`
}

// Status returns a summary of the download queue counters.
func (s *QueueService) Status(ctx context.Context) (*QueueStatusSummary, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QueueStatusSummary{}).
		Get("/api/v3/queue/status")
	if err != nil {
		return nil, fmt.Errorf("radarr: get queue status: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get queue status: %w", err)
	}

	return ptrResult[QueueStatusSummary](resp)
}

// ---------------------------------------------------------------------------
// Grab
// ---------------------------------------------------------------------------

// Grab triggers an immediate download for a queue item by ID.
func (s *QueueService) Grab(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Post("/api/v3/queue/grab/{id}")
	if err != nil {
		return fmt.Errorf("radarr: grab queue item %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: grab queue item %d: %w", id, err)
	}

	return nil
}

// GrabBulk triggers an immediate download for multiple queue items.
func (s *QueueService) GrabBulk(ctx context.Context, body QueueBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/queue/grab/bulk")
	if err != nil {
		return fmt.Errorf("radarr: grab queue items bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: grab queue items bulk: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// DeleteQueueOption is a functional option for QueueService.Delete and DeleteBulk.
type DeleteQueueOption func(*resty.Request)

// WithQueueRemoveFromClient also removes the download from the download client when v is true.
func WithQueueRemoveFromClient(v bool) DeleteQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("removeFromClient", strconv.FormatBool(v))
	}
}

// WithQueueBlocklist adds the release to the blocklist when v is true.
func WithQueueBlocklist(v bool) DeleteQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("blocklist", strconv.FormatBool(v))
	}
}

// WithQueueSkipRedownload skips the automatic re-download when v is true.
func WithQueueSkipRedownload(v bool) DeleteQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("skipRedownload", strconv.FormatBool(v))
	}
}

// WithQueueChangeCategory changes the download client category when v is true.
func WithQueueChangeCategory(v bool) DeleteQueueOption {
	return func(r *resty.Request) {
		r.SetQueryParam("changeCategory", strconv.FormatBool(v))
	}
}

// Delete removes a single queue item by ID.
func (s *QueueService) Delete(ctx context.Context, id int, opts ...DeleteQueueOption) error {
	req := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id))
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Delete("/api/v3/queue/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete queue item %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete queue item %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple queue items in a single request.
func (s *QueueService) DeleteBulk(ctx context.Context, body QueueBulkRequest, opts ...DeleteQueueOption) error {
	req := s.client.R().
		SetContext(ctx).
		SetBody(body)
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Delete("/api/v3/queue/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete queue items bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete queue items bulk: %w", err)
	}

	return nil
}
