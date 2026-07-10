package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// ManualImportService provides methods for the /manualimport endpoint.
type ManualImportService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ManualImportItem is a file candidate for manual import.
type ManualImportItem struct {
	ID                int                       `json:"id"`
	Path              *string                   `json:"path,omitempty"`
	RelativePath      *string                   `json:"relativePath,omitempty"`
	FolderName        *string                   `json:"folderName,omitempty"`
	Name              *string                   `json:"name,omitempty"`
	Size              int64                     `json:"size"`
	Movie             *Movie                    `json:"movie,omitempty"`
	MovieFileID       *int                      `json:"movieFileId,omitempty"`
	ReleaseGroup      *string                   `json:"releaseGroup,omitempty"`
	Quality           QualityModel              `json:"quality"`
	Languages         []Language                `json:"languages,omitempty"`
	QualityWeight     int                       `json:"qualityWeight"`
	DownloadID        *string                   `json:"downloadId,omitempty"`
	CustomFormats     []CustomFormat            `json:"customFormats,omitempty"`
	CustomFormatScore int                       `json:"customFormatScore"`
	IndexerFlags      int                       `json:"indexerFlags"`
	Rejections        []ImportRejectionResource `json:"rejections,omitempty"`
}

// ManualImportReprocessItem is the body element for a reprocess request.
type ManualImportReprocessItem struct {
	ID                int                       `json:"id"`
	Path              *string                   `json:"path,omitempty"`
	MovieID           int                       `json:"movieId"`
	Movie             *Movie                    `json:"movie,omitempty"`
	Quality           QualityModel              `json:"quality"`
	Languages         []Language                `json:"languages,omitempty"`
	ReleaseGroup      *string                   `json:"releaseGroup,omitempty"`
	DownloadID        *string                   `json:"downloadId,omitempty"`
	CustomFormats     []CustomFormat            `json:"customFormats,omitempty"`
	CustomFormatScore int                       `json:"customFormatScore"`
	IndexerFlags      int                       `json:"indexerFlags"`
	Rejections        []ImportRejectionResource `json:"rejections,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListManualImportOption is a functional option for ManualImportService.List.
type ListManualImportOption func(*resty.Request)

// WithManualImportFolder sets the folder path to scan for importable files.
func WithManualImportFolder(folder string) ListManualImportOption {
	return func(r *resty.Request) {
		r.SetQueryParam("folder", folder)
	}
}

// WithManualImportDownloadID filters results to a specific download client ID.
func WithManualImportDownloadID(id string) ListManualImportOption {
	return func(r *resty.Request) {
		r.SetQueryParam("downloadId", id)
	}
}

// WithManualImportMovieID filters results to a specific movie ID.
func WithManualImportMovieID(id int) ListManualImportOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// WithManualImportFilterExistingFiles excludes already-imported files when v is true.
func WithManualImportFilterExistingFiles(v bool) ListManualImportOption {
	return func(r *resty.Request) {
		r.SetQueryParam("filterExistingFiles", strconv.FormatBool(v))
	}
}

// List returns files available for manual import.
func (s *ManualImportService) List(ctx context.Context, opts ...ListManualImportOption) ([]ManualImportItem, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]ManualImportItem{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/manualimport")
	if err != nil {
		return nil, fmt.Errorf("radarr: list manual import: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list manual import: %w", err)
	}

	return derefResult[[]ManualImportItem](resp)
}

// ---------------------------------------------------------------------------
// Reprocess (PUT)
// ---------------------------------------------------------------------------

// Reprocess sends one or more import items for reprocessing.
func (s *ManualImportService) Reprocess(ctx context.Context, body []ManualImportReprocessItem) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Put("/api/v3/manualimport")
	if err != nil {
		return fmt.Errorf("radarr: reprocess manual import: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: reprocess manual import: %w", err)
	}

	return nil
}
