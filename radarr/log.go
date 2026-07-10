package radarr

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// LogService provides methods for the /log endpoint.
type LogService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// LogRecord is a single log entry from Radarr.
type LogRecord struct {
	ID            int       `json:"id"`
	Time          time.Time `json:"time"`
	Exception     *string   `json:"exception,omitempty"`
	ExceptionType *string   `json:"exceptionType,omitempty"`
	Level         *string   `json:"level,omitempty"`
	Logger        *string   `json:"logger,omitempty"`
	Message       *string   `json:"message,omitempty"`
	Method        *string   `json:"method,omitempty"`
}

// LogFile is a metadata descriptor for a log file on disk.
type LogFile struct {
	ID            int       `json:"id"`
	Filename      *string   `json:"filename,omitempty"`
	LastWriteTime time.Time `json:"lastWriteTime"`
	ContentsURL   *string   `json:"contentsUrl,omitempty"`
	DownloadURL   *string   `json:"downloadUrl,omitempty"`
}

// ---------------------------------------------------------------------------
// List (paged log records)
// ---------------------------------------------------------------------------

// ListLogOption is a functional option for LogService.List.
type ListLogOption func(*resty.Request)

// WithLogPage requests a specific page number.
func WithLogPage(page int) ListLogOption {
	return func(r *resty.Request) {
		r.SetQueryParam("page", strconv.Itoa(page))
	}
}

// WithLogPageSize sets the number of records per page.
func WithLogPageSize(n int) ListLogOption {
	return func(r *resty.Request) {
		r.SetQueryParam("pageSize", strconv.Itoa(n))
	}
}

// WithLogSortKey sets the field used to sort results.
func WithLogSortKey(key string) ListLogOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortKey", key)
	}
}

// WithLogSortDirection sets the sort order.
func WithLogSortDirection(dir SortDirection) ListLogOption {
	return func(r *resty.Request) {
		r.SetQueryParam("sortDirection", string(dir))
	}
}

// WithLogLevel filters log entries to a specific severity level.
func WithLogLevel(level string) ListLogOption {
	return func(r *resty.Request) {
		r.SetQueryParam("level", level)
	}
}

// List returns a single page of log records.
func (s *LogService) List(ctx context.Context, opts ...ListLogOption) (*PagedResult[LogRecord], error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&PagedResult[LogRecord]{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/log")
	if err != nil {
		return nil, fmt.Errorf("radarr: list log: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list log: %w", err)
	}

	return ptrResult[PagedResult[LogRecord]](resp)
}

// ListWithPagination iterates over all log pages automatically.
func (s *LogService) ListWithPagination(ctx context.Context, opts ...ListLogOption) iter.Seq2[PagedResult[LogRecord], error] {
	return paginate(ctx, s.List, opts...)
}

// ---------------------------------------------------------------------------
// Log files
// ---------------------------------------------------------------------------

// ListFiles returns the list of log files on disk.
func (s *LogService) ListFiles(ctx context.Context) ([]LogFile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]LogFile{}).
		Get("/api/v3/log/file")
	if err != nil {
		return nil, fmt.Errorf("radarr: list log files: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list log files: %w", err)
	}

	return derefResult[[]LogFile](resp)
}

// GetFile returns the content of a specific log file by filename.
func (s *LogService) GetFile(ctx context.Context, filename string) (string, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("filename", filename).
		Get("/api/v3/log/file/{filename}")
	if err != nil {
		return "", fmt.Errorf("radarr: get log file %s: %w", filename, err)
	}

	if err := checkResponse(resp); err != nil {
		return "", fmt.Errorf("radarr: get log file %s: %w", filename, err)
	}

	return string(resp.Body()), nil
}

// ListUpdateFiles returns the list of update log files on disk.
func (s *LogService) ListUpdateFiles(ctx context.Context) ([]LogFile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]LogFile{}).
		Get("/api/v3/log/file/update")
	if err != nil {
		return nil, fmt.Errorf("radarr: list update log files: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list update log files: %w", err)
	}

	return derefResult[[]LogFile](resp)
}

// GetUpdateFile returns the content of a specific update log file by filename.
func (s *LogService) GetUpdateFile(ctx context.Context, filename string) (string, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("filename", filename).
		Get("/api/v3/log/file/update/{filename}")
	if err != nil {
		return "", fmt.Errorf("radarr: get update log file %s: %w", filename, err)
	}

	if err := checkResponse(resp); err != nil {
		return "", fmt.Errorf("radarr: get update log file %s: %w", filename, err)
	}

	return string(resp.Body()), nil
}
