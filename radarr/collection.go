package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// CollectionService provides methods for the /collection endpoint.
type CollectionService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CollectionMovie is a movie that belongs to a collection but may not be in the library.
type CollectionMovie struct {
	TmdbID     int             `json:"tmdbId"`
	ImdbID     *string         `json:"imdbId,omitempty"`
	Title      *string         `json:"title,omitempty"`
	CleanTitle *string         `json:"cleanTitle,omitempty"`
	SortTitle  *string         `json:"sortTitle,omitempty"`
	Status     MovieStatusType `json:"status,omitempty"`
	Overview   *string         `json:"overview,omitempty"`
	Runtime    int             `json:"runtime"`
	Images     []MediaCover    `json:"images,omitempty"`
	Year       int             `json:"year"`
	Ratings    *Ratings        `json:"ratings,omitempty"`
	Genres     []string        `json:"genres,omitempty"`
	Folder     *string         `json:"folder,omitempty"`
	IsExisting bool            `json:"isExisting"`
	IsExcluded bool            `json:"isExcluded"`
}

// Collection groups related movies into a TMDB collection.
type Collection struct {
	ID                  int               `json:"id"`
	Title               *string           `json:"title,omitempty"`
	SortTitle           *string           `json:"sortTitle,omitempty"`
	TmdbID              int               `json:"tmdbId"`
	Images              []MediaCover      `json:"images,omitempty"`
	Overview            *string           `json:"overview,omitempty"`
	Monitored           bool              `json:"monitored"`
	RootFolderPath      *string           `json:"rootFolderPath,omitempty"`
	QualityProfileID    int               `json:"qualityProfileId"`
	SearchOnAdd         bool              `json:"searchOnAdd"`
	MinimumAvailability MovieStatusType   `json:"minimumAvailability,omitempty"`
	Movies              []CollectionMovie `json:"movies,omitempty"`
	MissingMovies       int               `json:"missingMovies"`
	Tags                []int             `json:"tags,omitempty"`
}

// CollectionUpdateRequest is the body for bulk-updating collections.
type CollectionUpdateRequest struct {
	CollectionIDs       []int           `json:"collectionIds,omitempty"`
	Monitored           *bool           `json:"monitored,omitempty"`
	MonitorMovies       *bool           `json:"monitorMovies,omitempty"`
	SearchOnAdd         *bool           `json:"searchOnAdd,omitempty"`
	QualityProfileID    *int            `json:"qualityProfileId,omitempty"`
	RootFolderPath      *string         `json:"rootFolderPath,omitempty"`
	MinimumAvailability MovieStatusType `json:"minimumAvailability,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all collections.
func (s *CollectionService) List(ctx context.Context) ([]Collection, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Collection{}).
		Get("/api/v3/collection")
	if err != nil {
		return nil, fmt.Errorf("radarr: list collections: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list collections: %w", err)
	}

	return derefResult[[]Collection](resp)
}

// Get returns a single collection by ID.
func (s *CollectionService) Get(ctx context.Context, id int) (*Collection, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Collection{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/collection/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get collection %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get collection %d: %w", id, err)
	}

	return ptrResult[Collection](resp)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// Update replaces a single collection by ID.
func (s *CollectionService) Update(ctx context.Context, id int, body Collection) (*Collection, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Collection{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/collection/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update collection %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update collection %d: %w", id, err)
	}

	return ptrResult[Collection](resp)
}

// UpdateBulk updates multiple collections at once.
func (s *CollectionService) UpdateBulk(ctx context.Context, body CollectionUpdateRequest) ([]Collection, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Collection{}).
		SetBody(body).
		Put("/api/v3/collection")
	if err != nil {
		return nil, fmt.Errorf("radarr: update collections bulk: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update collections bulk: %w", err)
	}

	return derefResult[[]Collection](resp)
}
