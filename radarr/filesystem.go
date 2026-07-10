package radarr

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

// FilesystemService provides methods for the /filesystem endpoint.
type FilesystemService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// FilesystemItem is a file or directory entry returned by Browse.
type FilesystemItem struct {
	Path         *string `json:"path,omitempty"`
	Name         *string `json:"name,omitempty"`
	LastModified *string `json:"lastModified,omitempty"`
	Size         *int64  `json:"size,omitempty"`
	Type         *string `json:"type,omitempty"`
}

// FilesystemResult holds directories and files from a Browse call.
type FilesystemResult struct {
	Parent      *string          `json:"parent,omitempty"`
	Directories []FilesystemItem `json:"directories,omitempty"`
	Files       []FilesystemItem `json:"files,omitempty"`
}

// ---------------------------------------------------------------------------
// Browse
// ---------------------------------------------------------------------------

// BrowseFilesystemOption is a functional option for FilesystemService.Browse.
type BrowseFilesystemOption func(*resty.Request)

// WithFilesystemPath sets the path to browse.
func WithFilesystemPath(path string) BrowseFilesystemOption {
	return func(r *resty.Request) {
		r.SetQueryParam("path", path)
	}
}

// WithFilesystemIncludeFiles includes files in the listing when v is true.
func WithFilesystemIncludeFiles(v bool) BrowseFilesystemOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeFiles", strconv.FormatBool(v))
	}
}

// WithFilesystemAllowFoldersWithoutTrailingSlashes includes folders without trailing slashes when v is true.
func WithFilesystemAllowFoldersWithoutTrailingSlashes(v bool) BrowseFilesystemOption {
	return func(r *resty.Request) {
		r.SetQueryParam("allowFoldersWithoutTrailingSlashes", strconv.FormatBool(v))
	}
}

// Browse lists the contents of a filesystem path.
func (s *FilesystemService) Browse(ctx context.Context, opts ...BrowseFilesystemOption) (*FilesystemResult, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&FilesystemResult{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/filesystem")
	if err != nil {
		return nil, fmt.Errorf("radarr: browse filesystem: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: browse filesystem: %w", err)
	}

	return ptrResult[FilesystemResult](resp)
}

// ---------------------------------------------------------------------------
// MediaFiles
// ---------------------------------------------------------------------------

// MediaFilesOption is a functional option for FilesystemService.MediaFiles.
type MediaFilesOption func(*resty.Request)

// WithMediaFilesPath sets the path to scan for media files.
func WithMediaFilesPath(path string) MediaFilesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("path", path)
	}
}

// MediaFiles returns media files found under the given path.
func (s *FilesystemService) MediaFiles(ctx context.Context, opts ...MediaFilesOption) ([]FilesystemItem, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]FilesystemItem{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/filesystem/mediafiles")
	if err != nil {
		return nil, fmt.Errorf("radarr: filesystem media files: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: filesystem media files: %w", err)
	}

	return derefResult[[]FilesystemItem](resp)
}

// ---------------------------------------------------------------------------
// Type
// ---------------------------------------------------------------------------

// FileTypeOption is a functional option for FilesystemService.FileType.
type FileTypeOption func(*resty.Request)

// WithFileTypePath sets the path to identify.
func WithFileTypePath(path string) FileTypeOption {
	return func(r *resty.Request) {
		r.SetQueryParam("path", path)
	}
}

// FileType returns the type ("folder" or "file") of the given path.
func (s *FilesystemService) FileType(ctx context.Context, opts ...FileTypeOption) (*string, error) {
	req := s.client.R().
		SetContext(ctx)
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/filesystem/type")
	if err != nil {
		return nil, fmt.Errorf("radarr: filesystem type: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: filesystem type: %w", err)
	}

	body := string(resp.Body())
	// Strip surrounding quotes if the API returns a JSON string.
	body = strings.Trim(body, `"`)

	return &body, nil
}
