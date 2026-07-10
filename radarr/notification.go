package radarr

import (
	"context"
	"fmt"
	"strconv"
)

// NotificationService provides methods for the /notification endpoint.
type NotificationService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// Notification is a configured notification connection.
type Notification struct {
	ID                                  int              `json:"id"`
	Name                                *string          `json:"name,omitempty"`
	Fields                              []Field          `json:"fields,omitempty"`
	ImplementationName                  *string          `json:"implementationName,omitempty"`
	Implementation                      *string          `json:"implementation,omitempty"`
	ConfigContract                      *string          `json:"configContract,omitempty"`
	InfoLink                            *string          `json:"infoLink,omitempty"`
	Message                             *ProviderMessage `json:"message,omitempty"`
	Tags                                []int            `json:"tags,omitempty"`
	Presets                             []Notification   `json:"presets,omitempty"`
	Link                                *string          `json:"link,omitempty"`
	OnGrab                              bool             `json:"onGrab"`
	OnDownload                          bool             `json:"onDownload"`
	OnUpgrade                           bool             `json:"onUpgrade"`
	OnRename                            bool             `json:"onRename"`
	OnMovieAdded                        bool             `json:"onMovieAdded"`
	OnMovieDelete                       bool             `json:"onMovieDelete"`
	OnMovieFileDelete                   bool             `json:"onMovieFileDelete"`
	OnMovieFileDeleteForUpgrade         bool             `json:"onMovieFileDeleteForUpgrade"`
	OnHealthIssue                       bool             `json:"onHealthIssue"`
	IncludeHealthWarnings               bool             `json:"includeHealthWarnings"`
	OnHealthRestored                    bool             `json:"onHealthRestored"`
	OnApplicationUpdate                 bool             `json:"onApplicationUpdate"`
	OnManualInteractionRequired         bool             `json:"onManualInteractionRequired"`
	SupportsOnGrab                      bool             `json:"supportsOnGrab"`
	SupportsOnDownload                  bool             `json:"supportsOnDownload"`
	SupportsOnUpgrade                   bool             `json:"supportsOnUpgrade"`
	SupportsOnRename                    bool             `json:"supportsOnRename"`
	SupportsOnMovieAdded                bool             `json:"supportsOnMovieAdded"`
	SupportsOnMovieDelete               bool             `json:"supportsOnMovieDelete"`
	SupportsOnMovieFileDelete           bool             `json:"supportsOnMovieFileDelete"`
	SupportsOnMovieFileDeleteForUpgrade bool             `json:"supportsOnMovieFileDeleteForUpgrade"`
	SupportsOnHealthIssue               bool             `json:"supportsOnHealthIssue"`
	SupportsOnHealthRestored            bool             `json:"supportsOnHealthRestored"`
	SupportsOnApplicationUpdate         bool             `json:"supportsOnApplicationUpdate"`
	SupportsOnManualInteractionRequired bool             `json:"supportsOnManualInteractionRequired"`
	TestCommand                         *string          `json:"testCommand,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get / Schema
// ---------------------------------------------------------------------------

// List returns all configured notifications.
func (s *NotificationService) List(ctx context.Context) ([]Notification, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Notification{}).
		Get("/api/v3/notification")
	if err != nil {
		return nil, fmt.Errorf("radarr: list notifications: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list notifications: %w", err)
	}

	return derefResult[[]Notification](resp)
}

// Get returns a single notification by ID.
func (s *NotificationService) Get(ctx context.Context, id int) (*Notification, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Notification{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/notification/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get notification %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get notification %d: %w", id, err)
	}

	return ptrResult[Notification](resp)
}

// Schema returns the available notification implementation schemas.
func (s *NotificationService) Schema(ctx context.Context) ([]Notification, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Notification{}).
		Get("/api/v3/notification/schema")
	if err != nil {
		return nil, fmt.Errorf("radarr: notification schema: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: notification schema: %w", err)
	}

	return derefResult[[]Notification](resp)
}

// ---------------------------------------------------------------------------
// Create / Update / Delete
// ---------------------------------------------------------------------------

// Create adds a new notification connection.
func (s *NotificationService) Create(ctx context.Context, body Notification) (*Notification, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Notification{}).
		SetBody(body).
		Post("/api/v3/notification")
	if err != nil {
		return nil, fmt.Errorf("radarr: create notification: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create notification: %w", err)
	}

	return ptrResult[Notification](resp)
}

// Update replaces a notification by ID.
func (s *NotificationService) Update(ctx context.Context, id int, body Notification) (*Notification, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Notification{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/notification/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update notification %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update notification %d: %w", id, err)
	}

	return ptrResult[Notification](resp)
}

// Delete removes a notification by ID.
func (s *NotificationService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/notification/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete notification %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete notification %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Test / Action
// ---------------------------------------------------------------------------

// Test validates a notification configuration.
func (s *NotificationService) Test(ctx context.Context, body Notification) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(body).
		Post("/api/v3/notification/test")
	if err != nil {
		return fmt.Errorf("radarr: test notification: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test notification: %w", err)
	}

	return nil
}

// TestAll validates all configured notifications.
func (s *NotificationService) TestAll(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/notification/testall")
	if err != nil {
		return fmt.Errorf("radarr: test all notifications: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: test all notifications: %w", err)
	}

	return nil
}

// Action executes a named action on a notification.
func (s *NotificationService) Action(ctx context.Context, name string, body Notification) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("name", name).
		SetBody(body).
		Post("/api/v3/notification/action/{name}")
	if err != nil {
		return fmt.Errorf("radarr: notification action %s: %w", name, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: notification action %s: %w", name, err)
	}

	return nil
}
