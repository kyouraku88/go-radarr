package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// QualityProfileService provides methods for the /qualityprofile endpoint.
type QualityProfileService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// QualityProfileItem is a quality tier entry within a profile.
type QualityProfileItem struct {
	ID      int                  `json:"id"`
	Name    *string              `json:"name,omitempty"`
	Quality *Quality             `json:"quality,omitempty"`
	Items   []QualityProfileItem `json:"items,omitempty"`
	Allowed bool                 `json:"allowed"`
}

// ProfileFormatItem associates a custom format with a score in a quality profile.
type ProfileFormatItem struct {
	ID     int     `json:"id"`
	Format int     `json:"format"`
	Name   *string `json:"name,omitempty"`
	Score  int     `json:"score"`
}

// QualityProfile defines which quality tiers Radarr will accept for a movie.
type QualityProfile struct {
	ID                    int                  `json:"id"`
	Name                  *string              `json:"name,omitempty"`
	UpgradeAllowed        bool                 `json:"upgradeAllowed"`
	Cutoff                int                  `json:"cutoff"`
	Items                 []QualityProfileItem `json:"items,omitempty"`
	MinFormatScore        int                  `json:"minFormatScore"`
	CutoffFormatScore     int                  `json:"cutoffFormatScore"`
	MinUpgradeFormatScore int                  `json:"minUpgradeFormatScore"`
	FormatItems           []ProfileFormatItem  `json:"formatItems,omitempty"`
	Language              Language             `json:"language"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all quality profiles.
func (s *QualityProfileService) List(ctx context.Context) ([]QualityProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]QualityProfile{}).
		Get("/api/v3/qualityprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: list quality profiles: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list quality profiles: %w", err)
	}

	return derefResult[[]QualityProfile](resp)
}

// Get returns a single quality profile by ID.
func (s *QualityProfileService) Get(ctx context.Context, id int) (*QualityProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/qualityprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get quality profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get quality profile %d: %w", id, err)
	}

	return ptrResult[QualityProfile](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new quality profile.
func (s *QualityProfileService) Create(ctx context.Context, body QualityProfile) (*QualityProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityProfile{}).
		SetBody(body).
		Post("/api/v3/qualityprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: create quality profile: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create quality profile: %w", err)
	}

	return ptrResult[QualityProfile](resp)
}

// Update replaces a quality profile by ID.
func (s *QualityProfileService) Update(ctx context.Context, id int, body QualityProfile) (*QualityProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&QualityProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/qualityprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update quality profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update quality profile %d: %w", id, err)
	}

	return ptrResult[QualityProfile](resp)
}

// Delete removes a quality profile by ID.
func (s *QualityProfileService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/qualityprofile/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete quality profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete quality profile %d: %w", id, err)
	}

	return nil
}
