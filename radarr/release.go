package radarr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// ReleaseService provides methods for the /release endpoint.
type ReleaseService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Release is a search result or manually pushed release.
type Release struct {
	ID                  int              `json:"id"`
	GUID                *string          `json:"guid,omitempty"`
	Quality             QualityModel     `json:"quality"`
	CustomFormats       []CustomFormat   `json:"customFormats,omitempty"`
	CustomFormatScore   int              `json:"customFormatScore"`
	QualityWeight       int              `json:"qualityWeight"`
	Age                 int              `json:"age"`
	AgeHours            float64          `json:"ageHours"`
	AgeMinutes          float64          `json:"ageMinutes"`
	Size                int64            `json:"size"`
	IndexerID           int              `json:"indexerId"`
	Indexer             *string          `json:"indexer,omitempty"`
	ReleaseGroup        *string          `json:"releaseGroup,omitempty"`
	SubGroup            *string          `json:"subGroup,omitempty"`
	ReleaseHash         *string          `json:"releaseHash,omitempty"`
	Title               *string          `json:"title,omitempty"`
	SceneSource         bool             `json:"sceneSource"`
	MovieTitles         []string         `json:"movieTitles,omitempty"`
	Languages           []Language       `json:"languages,omitempty"`
	MappedMovieID       *int             `json:"mappedMovieId,omitempty"`
	Approved            bool             `json:"approved"`
	TemporarilyRejected bool             `json:"temporarilyRejected"`
	Rejected            bool             `json:"rejected"`
	TmdbID              int              `json:"tmdbId"`
	ImdbID              int              `json:"imdbId"`
	Rejections          []string         `json:"rejections,omitempty"`
	PublishDate         time.Time        `json:"publishDate"`
	CommentURL          *string          `json:"commentUrl,omitempty"`
	DownloadURL         *string          `json:"downloadUrl,omitempty"`
	InfoURL             *string          `json:"infoUrl,omitempty"`
	MovieRequested      bool             `json:"movieRequested"`
	DownloadAllowed     bool             `json:"downloadAllowed"`
	ReleaseWeight       int              `json:"releaseWeight"`
	Edition             *string          `json:"edition,omitempty"`
	MagnetURL           *string          `json:"magnetUrl,omitempty"`
	InfoHash            *string          `json:"infoHash,omitempty"`
	Seeders             *int             `json:"seeders,omitempty"`
	Leechers            *int             `json:"leechers,omitempty"`
	Protocol            DownloadProtocol `json:"protocol,omitempty"`
	IndexerFlags        []string         `json:"indexerFlags,omitempty"`
	MovieID             *int             `json:"movieId,omitempty"`
	DownloadClientID    *int             `json:"downloadClientId,omitempty"`
	DownloadClient      *string          `json:"downloadClient,omitempty"`
	ShouldOverride      *bool            `json:"shouldOverride,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListReleasesOption is a functional option for ReleaseService.List.
type ListReleasesOption func(*resty.Request)

// WithReleaseMovieID filters search results to a specific movie ID.
func WithReleaseMovieID(id int) ListReleasesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// List searches for releases for a given movie.
func (s *ReleaseService) List(ctx context.Context, opts ...ListReleasesOption) ([]Release, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Release{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/release")
	if err != nil {
		return nil, fmt.Errorf("radarr: list releases: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list releases: %w", err)
	}

	return derefResult[[]Release](resp)
}

// ---------------------------------------------------------------------------
// Download / Push
// ---------------------------------------------------------------------------

// Download adds a release to the download queue.
func (s *ReleaseService) Download(ctx context.Context, body Release) (*Release, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Release{}).
		SetBody(body).
		Post("/api/v3/release")
	if err != nil {
		return nil, fmt.Errorf("radarr: download release: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: download release: %w", err)
	}

	return ptrResult[Release](resp)
}

// Push sends a release directly to the download client without searching.
func (s *ReleaseService) Push(ctx context.Context, body Release) ([]Release, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Release{}).
		SetBody(body).
		Post("/api/v3/release/push")
	if err != nil {
		return nil, fmt.Errorf("radarr: push release: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: push release: %w", err)
	}

	return derefResult[[]Release](resp)
}
