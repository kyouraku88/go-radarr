package radarr

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ParseService provides methods for the /parse endpoint.
type ParseService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ParsedMovieInfo is the result of parsing a release title.
type ParsedMovieInfo struct {
	MovieTitles        []string     `json:"movieTitles,omitempty"`
	OriginalTitle      *string      `json:"originalTitle,omitempty"`
	ReleaseTitle       *string      `json:"releaseTitle,omitempty"`
	SimpleReleaseTitle *string      `json:"simpleReleaseTitle,omitempty"`
	Quality            QualityModel `json:"quality"`
	Languages          []Language   `json:"languages,omitempty"`
	ReleaseGroup       *string      `json:"releaseGroup,omitempty"`
	ReleaseHash        *string      `json:"releaseHash,omitempty"`
	Edition            *string      `json:"edition,omitempty"`
	Year               int          `json:"year"`
	ImdbID             *string      `json:"imdbId,omitempty"`
	TmdbID             int          `json:"tmdbId"`
	HardcodedSubs      *string      `json:"hardcodedSubs,omitempty"`
	MovieTitle         *string      `json:"movieTitle,omitempty"`
	PrimaryMovieTitle  *string      `json:"primaryMovieTitle,omitempty"`
}

// ParseResult is the full response from the /parse endpoint.
type ParseResult struct {
	ID                int              `json:"id"`
	Title             *string          `json:"title,omitempty"`
	ParsedMovieInfo   *ParsedMovieInfo `json:"parsedMovieInfo,omitempty"`
	Movie             *Movie           `json:"movie,omitempty"`
	Languages         []Language       `json:"languages,omitempty"`
	CustomFormats     []CustomFormat   `json:"customFormats,omitempty"`
	CustomFormatScore int              `json:"customFormatScore"`
}

// ---------------------------------------------------------------------------
// Parse
// ---------------------------------------------------------------------------

// ParseOption is a functional option for ParseService.Parse.
type ParseOption func(*resty.Request)

// WithParseTitle sets the release title to parse.
func WithParseTitle(title string) ParseOption {
	return func(r *resty.Request) {
		r.SetQueryParam("title", title)
	}
}

// Parse parses a release title and returns matched movie information.
func (s *ParseService) Parse(ctx context.Context, opts ...ParseOption) (*ParseResult, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&ParseResult{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/parse")
	if err != nil {
		return nil, fmt.Errorf("radarr: parse: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: parse: %w", err)
	}

	return ptrResult[ParseResult](resp)
}
