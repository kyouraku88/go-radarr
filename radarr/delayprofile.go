package radarr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// DelayProfileService provides methods for the /delayprofile endpoint.
type DelayProfileService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// DelayProfile controls the delay before Radarr downloads a release.
type DelayProfile struct {
	ID                             int              `json:"id"`
	EnableUsenet                   bool             `json:"enableUsenet"`
	EnableTorrent                  bool             `json:"enableTorrent"`
	PreferredProtocol              DownloadProtocol `json:"preferredProtocol,omitempty"`
	UsenetDelay                    int              `json:"usenetDelay"`
	TorrentDelay                   int              `json:"torrentDelay"`
	BypassIfHighestQuality         bool             `json:"bypassIfHighestQuality"`
	BypassIfAboveCustomFormatScore bool             `json:"bypassIfAboveCustomFormatScore"`
	MinimumCustomFormatScore       int              `json:"minimumCustomFormatScore"`
	Order                          int              `json:"order"`
	Tags                           []int            `json:"tags,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all delay profiles.
func (s *DelayProfileService) List(ctx context.Context) ([]DelayProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]DelayProfile{}).
		Get("/api/v3/delayprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: list delay profiles: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list delay profiles: %w", err)
	}

	return derefResult[[]DelayProfile](resp)
}

// Get returns a single delay profile by ID.
func (s *DelayProfileService) Get(ctx context.Context, id int) (*DelayProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DelayProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/delayprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get delay profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get delay profile %d: %w", id, err)
	}

	return ptrResult[DelayProfile](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete / Reorder
// ---------------------------------------------------------------------------

// Create adds a new delay profile.
func (s *DelayProfileService) Create(ctx context.Context, body DelayProfile) (*DelayProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DelayProfile{}).
		SetBody(body).
		Post("/api/v3/delayprofile")
	if err != nil {
		return nil, fmt.Errorf("radarr: create delay profile: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create delay profile: %w", err)
	}

	return ptrResult[DelayProfile](resp)
}

// Update replaces a delay profile by ID.
func (s *DelayProfileService) Update(ctx context.Context, id int, body DelayProfile) (*DelayProfile, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DelayProfile{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/delayprofile/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update delay profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update delay profile %d: %w", id, err)
	}

	return ptrResult[DelayProfile](resp)
}

// Delete removes a delay profile by ID.
func (s *DelayProfileService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/delayprofile/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete delay profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete delay profile %d: %w", id, err)
	}

	return nil
}

// ReorderOption is a functional option for DelayProfileService.Reorder.
type ReorderOption func(*resty.Request)

// WithReorderAfter sets the ID of the delay profile that the moved profile should appear after.
func WithReorderAfter(afterID int) ReorderOption {
	return func(r *resty.Request) {
		r.SetQueryParam("after", strconv.Itoa(afterID))
	}
}

// Reorder moves a delay profile to a new position in the ordering.
func (s *DelayProfileService) Reorder(ctx context.Context, id int, opts ...ReorderOption) error {
	req := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id))
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Put("/api/v3/delayprofile/reorder/{id}")
	if err != nil {
		return fmt.Errorf("radarr: reorder delay profile %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: reorder delay profile %d: %w", id, err)
	}

	return nil
}
