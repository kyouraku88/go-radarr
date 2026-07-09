package radarr

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// CommandService provides methods for the /command endpoint.
type CommandService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CommandPriority controls the scheduling priority of a command.
type CommandPriority string

// Command priority values.
const (
	CommandPriorityNormal CommandPriority = "normal"
	CommandPriorityHigh   CommandPriority = "high"
	CommandPriorityLow    CommandPriority = "low"
)

// CommandStatus reports the current execution state of a command.
type CommandStatus string

// Command status values.
const (
	CommandStatusQueued    CommandStatus = "queued"
	CommandStatusStarted   CommandStatus = "started"
	CommandStatusCompleted CommandStatus = "completed"
	CommandStatusFailed    CommandStatus = "failed"
	CommandStatusAborted   CommandStatus = "aborted"
	CommandStatusCancelled CommandStatus = "cancelled"
	CommandStatusOrphaned  CommandStatus = "orphaned"
)

// CommandResult indicates whether a completed command succeeded.
type CommandResult string

// Command result values.
const (
	CommandResultUnknown      CommandResult = "unknown"
	CommandResultSuccessful   CommandResult = "successful"
	CommandResultUnsuccessful CommandResult = "unsuccessful"
)

// CommandTrigger identifies what initiated a command.
type CommandTrigger string

// Command trigger values.
const (
	CommandTriggerUnspecified CommandTrigger = "unspecified"
	CommandTriggerManual      CommandTrigger = "manual"
	CommandTriggerScheduled   CommandTrigger = "scheduled"
)

// CommandBody is the payload embedded inside a CommandRecord.
type CommandBody struct {
	SendUpdatesToClient bool           `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool           `json:"updateScheduledTask,omitempty"`
	CompletionMessage   *string        `json:"completionMessage,omitempty"`
	RequiresDiskAccess  bool           `json:"requiresDiskAccess,omitempty"`
	IsExclusive         bool           `json:"isExclusive,omitempty"`
	IsTypeExclusive     bool           `json:"isTypeExclusive,omitempty"`
	IsLongRunning       bool           `json:"isLongRunning,omitempty"`
	Name                *string        `json:"name,omitempty"`
	LastExecutionTime   *time.Time     `json:"lastExecutionTime,omitempty"`
	LastStartTime       *time.Time     `json:"lastStartTime,omitempty"`
	Trigger             CommandTrigger `json:"trigger,omitempty"`
	SuppressMessages    bool           `json:"suppressMessages,omitempty"`
	ClientUserAgent     *string        `json:"clientUserAgent,omitempty"`
}

// CommandRecord represents a queued or completed Radarr command.
type CommandRecord struct {
	ID                  int             `json:"id"`
	Name                *string         `json:"name,omitempty"`
	CommandName         *string         `json:"commandName,omitempty"`
	Message             *string         `json:"message,omitempty"`
	Body                *CommandBody    `json:"body,omitempty"`
	Priority            CommandPriority `json:"priority,omitempty"`
	Status              CommandStatus   `json:"status,omitempty"`
	Result              CommandResult   `json:"result,omitempty"`
	Queued              time.Time       `json:"queued"`
	Started             *time.Time      `json:"started,omitempty"`
	Ended               *time.Time      `json:"ended,omitempty"`
	Duration            *string         `json:"duration,omitempty"`
	Exception           *string         `json:"exception,omitempty"`
	Trigger             CommandTrigger  `json:"trigger,omitempty"`
	ClientUserAgent     *string         `json:"clientUserAgent,omitempty"`
	StateChangeTime     *time.Time      `json:"stateChangeTime,omitempty"`
	SendUpdatesToClient bool            `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool            `json:"updateScheduledTask"`
	LastExecutionTime   *time.Time      `json:"lastExecutionTime,omitempty"`
}

// ---------------------------------------------------------------------------
// List / Get
// ---------------------------------------------------------------------------

// List returns all commands in the queue.
func (s *CommandService) List(ctx context.Context) ([]CommandRecord, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]CommandRecord{}).
		Get("/api/v3/command")
	if err != nil {
		return nil, fmt.Errorf("radarr: list commands: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list commands: %w", err)
	}

	return derefResult[[]CommandRecord](resp)
}

// Get returns a single command by ID.
func (s *CommandService) Get(ctx context.Context, id int) (*CommandRecord, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CommandRecord{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/command/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get command %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get command %d: %w", id, err)
	}

	return ptrResult[CommandRecord](resp)
}

// ---------------------------------------------------------------------------
// Create / Delete
// ---------------------------------------------------------------------------

// Create dispatches a command to Radarr. Populate at minimum CommandRecord.Name
// (e.g. "RefreshMovie", "RescanMovie", "MissingMoviesSearch").
func (s *CommandService) Create(ctx context.Context, body CommandRecord) (*CommandRecord, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&CommandRecord{}).
		SetBody(body).
		Post("/api/v3/command")
	if err != nil {
		return nil, fmt.Errorf("radarr: create command: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: create command: %w", err)
	}

	return ptrResult[CommandRecord](resp)
}

// Delete removes a command from the queue by ID.
func (s *CommandService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/command/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete command %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete command %d: %w", id, err)
	}

	return nil
}
