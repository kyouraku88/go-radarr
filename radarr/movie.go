package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// MovieService provides methods for the /movie endpoint.
type MovieService service

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// ListMoviesOption is a functional option for MovieService.List.
type ListMoviesOption func(*resty.Request)

// WithTmdbID filters results to a specific TMDB movie ID.
func WithTmdbID(id int) ListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("tmdbId", strconv.Itoa(id))
	}
}

// WithLanguageID filters results by language ID.
func WithLanguageID(id int) ListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("languageId", strconv.Itoa(id))
	}
}

// WithExcludeLocalCovers excludes locally stored cover art from results when v is true.
func WithExcludeLocalCovers(v bool) ListMoviesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("excludeLocalCovers", strconv.FormatBool(v))
	}
}

// List returns all movies in the library.
func (s *MovieService) List(ctx context.Context, opts ...ListMoviesOption) ([]Movie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/movie")
	if err != nil {
		return nil, fmt.Errorf("radarr: list movies: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list movies: %w", err)
	}

	return derefResult[[]Movie](resp)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

// Get returns a single movie by ID.
func (s *MovieService) Get(ctx context.Context, id int) (*Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Movie{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/movie/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get movie %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get movie %d: %w", id, err)
	}

	return ptrResult[Movie](resp)
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

// Create adds a movie to the Radarr library.
func (s *MovieService) Create(ctx context.Context, body Movie) (*Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Movie{}).
		SetBody(body).
		Post("/api/v3/movie")
	if err != nil {
		return nil, fmt.Errorf("radarr: create movie: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create movie: %w", err)
	}

	return ptrResult[Movie](resp)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// UpdateMovieOption is a functional option for MovieService.Update.
type UpdateMovieOption func(*resty.Request)

// WithMoveFiles moves the movie files on disk when v is true.
func WithMoveFiles(v bool) UpdateMovieOption {
	return func(r *resty.Request) {
		r.SetQueryParam("moveFiles", strconv.FormatBool(v))
	}
}

// Update replaces a movie's metadata by ID.
func (s *MovieService) Update(ctx context.Context, id int, body Movie, opts ...UpdateMovieOption) (*Movie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&Movie{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body)
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Put("/api/v3/movie/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update movie %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update movie %d: %w", id, err)
	}

	return ptrResult[Movie](resp)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// DeleteMovieOption is a functional option for MovieService.Delete.
type DeleteMovieOption func(*resty.Request)

// WithDeleteFiles also deletes the movie files on disk when v is true.
func WithDeleteFiles(v bool) DeleteMovieOption {
	return func(r *resty.Request) {
		r.SetQueryParam("deleteFiles", strconv.FormatBool(v))
	}
}

// WithAddImportExclusion adds the movie to the import exclusion list when v is true.
func WithAddImportExclusion(v bool) DeleteMovieOption {
	return func(r *resty.Request) {
		r.SetQueryParam("addImportExclusion", strconv.FormatBool(v))
	}
}

// Delete removes a movie from the library by ID.
func (s *MovieService) Delete(ctx context.Context, id int, opts ...DeleteMovieOption) error {
	req := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id))
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Delete("/api/v3/movie/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete movie %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete movie %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Lookup
// ---------------------------------------------------------------------------

// LookupMovieOption is a functional option for MovieService.Lookup.
type LookupMovieOption func(*resty.Request)

// WithLookupTerm sets the search term for a movie lookup.
func WithLookupTerm(term string) LookupMovieOption {
	return func(r *resty.Request) {
		r.SetQueryParam("term", term)
	}
}

// Lookup searches for movies by term (title, IMDB ID, TMDB ID, etc.).
func (s *MovieService) Lookup(ctx context.Context, opts ...LookupMovieOption) ([]Movie, error) {
	req := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{})
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/movie/lookup")
	if err != nil {
		return nil, fmt.Errorf("radarr: movie lookup: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: movie lookup: %w", err)
	}

	return derefResult[[]Movie](resp)
}

// LookupByIMDB looks up a movie by its IMDB ID.
func (s *MovieService) LookupByIMDB(ctx context.Context, imdbID string) ([]Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{}).
		SetQueryParam("imdbId", imdbID).
		Get("/api/v3/movie/lookup/imdb")
	if err != nil {
		return nil, fmt.Errorf("radarr: movie lookup by imdb %s: %w", imdbID, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: movie lookup by imdb %s: %w", imdbID, err)
	}

	return derefResult[[]Movie](resp)
}

// LookupByTMDB looks up a movie by its TMDB ID.
func (s *MovieService) LookupByTMDB(ctx context.Context, tmdbID int) ([]Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{}).
		SetQueryParam("tmdbId", strconv.Itoa(tmdbID)).
		Get("/api/v3/movie/lookup/tmdb")
	if err != nil {
		return nil, fmt.Errorf("radarr: movie lookup by tmdb %d: %w", tmdbID, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: movie lookup by tmdb %d: %w", tmdbID, err)
	}

	return derefResult[[]Movie](resp)
}

// ---------------------------------------------------------------------------
// Import
// ---------------------------------------------------------------------------

// Import adds multiple movies to the library in bulk.
func (s *MovieService) Import(ctx context.Context, body []Movie) ([]Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Movie{}).
		SetBody(body).
		Post("/api/v3/movie/import")
	if err != nil {
		return nil, fmt.Errorf("radarr: import movies: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: import movies: %w", err)
	}

	return derefResult[[]Movie](resp)
}

// ---------------------------------------------------------------------------
// GetFolder
// ---------------------------------------------------------------------------

// GetFolder returns the folder information for a movie by ID.
func (s *MovieService) GetFolder(ctx context.Context, id int) (*Movie, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Movie{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/movie/{id}/folder")
	if err != nil {
		return nil, fmt.Errorf("radarr: get movie folder %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get movie folder %d: %w", id, err)
	}

	return ptrResult[Movie](resp)
}
