package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// ExtraFileService provides methods for the /extrafile endpoint.
type ExtraFileService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ExtraFile is an extra file (e.g. subtitle) associated with a movie file.
type ExtraFile struct {
	ID           int           `json:"id"`
	MovieID      int           `json:"movieId"`
	MovieFileID  *int          `json:"movieFileId,omitempty"`
	RelativePath *string       `json:"relativePath,omitempty"`
	Extension    *string       `json:"extension,omitempty"`
	LanguageTags []string      `json:"languageTags,omitempty"`
	Title        *string       `json:"title,omitempty"`
	Type         ExtraFileType `json:"type,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListExtraFilesOption is a functional option for ExtraFileService.List.
type ListExtraFilesOption func(*resty.Request)

// WithExtraFileMovieID filters extra files to a specific movie ID.
func WithExtraFileMovieID(id int) ListExtraFilesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieId", strconv.Itoa(id))
	}
}

// List returns extra files, optionally filtered by movie ID.
func (s *ExtraFileService) List(ctx context.Context, opts ...ListExtraFilesOption) ([]ExtraFile, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]ExtraFile{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/extrafile")
	if err != nil {
		return nil, fmt.Errorf("radarr: list extra files: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list extra files: %w", err)
	}

	return derefResult[[]ExtraFile](resp)
}
