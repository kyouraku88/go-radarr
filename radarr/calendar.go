package radarr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// CalendarService provides methods for the /calendar endpoint.
type CalendarService service

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListCalendarOption is a functional option for CalendarService.List.
type ListCalendarOption func(*resty.Request)

// WithCalendarStart filters results to movies released on or after t.
func WithCalendarStart(t time.Time) ListCalendarOption {
	return func(r *resty.Request) {
		r.SetQueryParam("start", t.UTC().Format(time.RFC3339))
	}
}

// WithCalendarEnd filters results to movies released before t.
func WithCalendarEnd(t time.Time) ListCalendarOption {
	return func(r *resty.Request) {
		r.SetQueryParam("end", t.UTC().Format(time.RFC3339))
	}
}

// WithCalendarUnmonitored includes unmonitored movies when v is true.
func WithCalendarUnmonitored(v bool) ListCalendarOption {
	return func(r *resty.Request) {
		r.SetQueryParam("unmonitored", strconv.FormatBool(v))
	}
}

// WithCalendarTags filters results to movies with the specified tag labels (comma-separated).
func WithCalendarTags(tags string) ListCalendarOption {
	return func(r *resty.Request) {
		r.SetQueryParam("tags", tags)
	}
}

// List returns movies with a release date within the requested window.
func (s *CalendarService) List(ctx context.Context, opts ...ListCalendarOption) ([]Movie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/calendar")
	if err != nil {
		return nil, fmt.Errorf("radarr: list calendar: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list calendar: %w", err)
	}

	return derefResult[[]Movie](resp)
}
