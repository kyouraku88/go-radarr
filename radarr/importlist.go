package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// ImportListService provides methods for the /importlist endpoint.
type ImportListService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ImportList is a configured import list source.
type ImportList struct {
	ID                  int              `json:"id"`
	Name                *string          `json:"name,omitempty"`
	Fields              []Field          `json:"fields,omitempty"`
	ImplementationName  *string          `json:"implementationName,omitempty"`
	Implementation      *string          `json:"implementation,omitempty"`
	ConfigContract      *string          `json:"configContract,omitempty"`
	InfoLink            *string          `json:"infoLink,omitempty"`
	Message             *ProviderMessage `json:"message,omitempty"`
	Tags                []int            `json:"tags,omitempty"`
	Presets             []ImportList     `json:"presets,omitempty"`
	Enabled             bool             `json:"enabled"`
	EnableAuto          bool             `json:"enableAuto"`
	Monitor             MonitorTypes     `json:"monitor,omitempty"`
	RootFolderPath      *string          `json:"rootFolderPath,omitempty"`
	QualityProfileID    int              `json:"qualityProfileId"`
	SearchOnAdd         bool             `json:"searchOnAdd"`
	MinimumAvailability MovieStatusType  `json:"minimumAvailability,omitempty"`
	ListType            ImportListType   `json:"listType,omitempty"`
	ListOrder           int              `json:"listOrder"`
	MinRefreshInterval  *string          `json:"minRefreshInterval,omitempty"`
}

// ImportListBulkRequest is the body for bulk update and delete operations on import lists.
type ImportListBulkRequest struct {
	IDs                 []int           `json:"ids,omitempty"`
	Tags                []int           `json:"tags,omitempty"`
	ApplyTags           ApplyTags       `json:"applyTags,omitempty"`
	Enabled             *bool           `json:"enabled,omitempty"`
	EnableAuto          *bool           `json:"enableAuto,omitempty"`
	RootFolderPath      *string         `json:"rootFolderPath,omitempty"`
	QualityProfileID    *int            `json:"qualityProfileId,omitempty"`
	MinimumAvailability MovieStatusType `json:"minimumAvailability,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get / Schema
// ---------------------------------------------------------------------------

// List returns all configured import lists.
func (s *ImportListService) List(ctx context.Context) ([]ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ImportList{}).
		Get("/api/v3/importlist")
	if err != nil {
		return nil, fmt.Errorf("radarr: list import lists: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list import lists: %w", err)
	}

	return derefResult[[]ImportList](resp)
}

// Get returns a single import list by ID.
func (s *ImportListService) Get(ctx context.Context, id int) (*ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportList{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/importlist/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get import list %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get import list %d: %w", id, err)
	}

	return ptrResult[ImportList](resp)
}

// Schema returns the available import list implementation schemas.
func (s *ImportListService) Schema(ctx context.Context) ([]ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ImportList{}).
		Get("/api/v3/importlist/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: import list schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: import list schema: %w", err)
	}

	return derefResult[[]ImportList](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new import list.
func (s *ImportListService) Create(ctx context.Context, body ImportList) (*ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportList{}).
		SetBody(body).
		Post("/api/v3/importlist")
	if err != nil {
		return nil, fmt.Errorf("radarr: create import list: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create import list: %w", err)
	}

	return ptrResult[ImportList](resp)
}

// Update replaces an import list by ID.
func (s *ImportListService) Update(ctx context.Context, id int, body ImportList) (*ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportList{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/importlist/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update import list %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update import list %d: %w", id, err)
	}

	return ptrResult[ImportList](resp)
}

// UpdateBulk updates multiple import lists at once.
func (s *ImportListService) UpdateBulk(ctx context.Context, body ImportListBulkRequest) ([]ImportList, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]ImportList{}).
		SetBody(body).
		Put("/api/v3/importlist/bulk")
	if err != nil {
		return nil, fmt.Errorf("radarr: update import lists bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update import lists bulk: %w", err)
	}

	return derefResult[[]ImportList](resp)
}

// Delete removes an import list by ID.
func (s *ImportListService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/importlist/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete import list %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete import list %d: %w", id, err)
	}

	return nil
}

// DeleteBulk removes multiple import lists in a single request.
func (s *ImportListService) DeleteBulk(ctx context.Context, body ImportListBulkRequest) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Delete("/api/v3/importlist/bulk")
	if err != nil {
		return fmt.Errorf("radarr: delete import lists bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete import lists bulk: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Test / Action / Movies
// ---------------------------------------------------------------------------

// Test validates an import list configuration.
func (s *ImportListService) Test(ctx context.Context, body ImportList) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/importlist/test")
	if err != nil {
		return fmt.Errorf("radarr: test import list: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test import list: %w", err)
	}

	return nil
}

// TestAll validates all configured import lists.
func (s *ImportListService) TestAll(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/importlist/testall")
	if err != nil {
		return fmt.Errorf("radarr: test all import lists: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test all import lists: %w", err)
	}

	return nil
}

// Action executes a named action on an import list.
func (s *ImportListService) Action(ctx context.Context, name string, body ImportList) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("name", name).
		SetBody(body).
		Post("/api/v3/importlist/action/{name}")
	if err != nil {
		return fmt.Errorf("radarr: import list action %s: %w", name, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: import list action %s: %w", name, err)
	}

	return nil
}

// ImportListMoviesOption is a functional option for ImportListService.Movies.
type ImportListMoviesOption func(*resty.Request)

// WithImportListMoviesIncludeRecommendations includes recommended movies when v is true.
func WithImportListMoviesIncludeRecommendations(v bool) ImportListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeRecommendations", strconv.FormatBool(v))
	}
}

// WithImportListMoviesIncludeTrending includes trending movies when v is true.
func WithImportListMoviesIncludeTrending(v bool) ImportListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includeTrending", strconv.FormatBool(v))
	}
}

// WithImportListMoviesIncludePopular includes popular movies when v is true.
func WithImportListMoviesIncludePopular(v bool) ImportListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("includePopular", strconv.FormatBool(v))
	}
}

// Movies returns the movies discovered by the import lists.
func (s *ImportListService) Movies(ctx context.Context, opts ...ImportListMoviesOption) ([]Movie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/importlist/movie")
	if err != nil {
		return nil, fmt.Errorf("radarr: import list movies: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: import list movies: %w", err)
	}

	return derefResult[[]Movie](resp)
}
