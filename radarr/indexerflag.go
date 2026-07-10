package radarr

import (
	"context"
	"fmt"
)

// IndexerFlagService provides methods for the /indexerflag endpoint.
type IndexerFlagService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// IndexerFlag is a flag that can be applied to an indexer result.
type IndexerFlag struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	NameLower *string `json:"nameLower,omitempty"`
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// List returns all indexer flags.
func (s *IndexerFlagService) List(ctx context.Context) ([]IndexerFlag, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]IndexerFlag{}).
		Get("/api/v3/indexerflag")
	if err != nil {
		return nil, fmt.Errorf("radarr: list indexer flags: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list indexer flags: %w", err)
	}

	return derefResult[[]IndexerFlag](resp)
}
