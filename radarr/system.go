package radarr

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// SystemService provides methods for the /system endpoint.
type SystemService service

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// RuntimeMode describes how the Radarr process is hosted.
type RuntimeMode string

// Runtime mode values.
const (
	RuntimeModeConsole RuntimeMode = "console"
	RuntimeModeService RuntimeMode = "service"
	RuntimeModeTray    RuntimeMode = "tray"
)

// DatabaseType identifies the backing database engine.
type DatabaseType string

// Database type values.
const (
	DatabaseTypeSQLite     DatabaseType = "sqLite"
	DatabaseTypePostgreSQL DatabaseType = "postgreSQL"
)

// AuthenticationType identifies the authentication mechanism in use.
type AuthenticationType string

// Authentication type values.
const (
	AuthenticationTypeNone     AuthenticationType = "none"
	AuthenticationTypeBasic    AuthenticationType = "basic"
	AuthenticationTypeForms    AuthenticationType = "forms"
	AuthenticationTypeExternal AuthenticationType = "external"
)

// UpdateMechanism describes how Radarr receives application updates.
type UpdateMechanism string

// Update mechanism values.
const (
	UpdateMechanismBuiltIn  UpdateMechanism = "builtIn"
	UpdateMechanismScript   UpdateMechanism = "script"
	UpdateMechanismExternal UpdateMechanism = "external"
	UpdateMechanismApt      UpdateMechanism = "apt"
	UpdateMechanismDocker   UpdateMechanism = "docker"
)

// SystemStatus holds runtime and environment information about the Radarr instance.
type SystemStatus struct {
	AppName                       *string            `json:"appName,omitempty"`
	InstanceName                  *string            `json:"instanceName,omitempty"`
	Version                       *string            `json:"version,omitempty"`
	BuildTime                     time.Time          `json:"buildTime"`
	IsDebug                       bool               `json:"isDebug"`
	IsProduction                  bool               `json:"isProduction"`
	IsAdmin                       bool               `json:"isAdmin"`
	IsUserInteractive             bool               `json:"isUserInteractive"`
	StartupPath                   *string            `json:"startupPath,omitempty"`
	AppData                       *string            `json:"appData,omitempty"`
	OsName                        *string            `json:"osName,omitempty"`
	OsVersion                     *string            `json:"osVersion,omitempty"`
	IsNetCore                     bool               `json:"isNetCore"`
	IsLinux                       bool               `json:"isLinux"`
	IsOsx                         bool               `json:"isOsx"`
	IsWindows                     bool               `json:"isWindows"`
	IsDocker                      bool               `json:"isDocker"`
	Mode                          RuntimeMode        `json:"mode,omitempty"`
	Branch                        *string            `json:"branch,omitempty"`
	DatabaseType                  DatabaseType       `json:"databaseType,omitempty"`
	DatabaseVersion               *string            `json:"databaseVersion,omitempty"`
	Authentication                AuthenticationType `json:"authentication,omitempty"`
	MigrationVersion              int                `json:"migrationVersion"`
	URLBase                       *string            `json:"urlBase,omitempty"`
	RuntimeVersion                *string            `json:"runtimeVersion,omitempty"`
	RuntimeName                   *string            `json:"runtimeName,omitempty"`
	StartTime                     time.Time          `json:"startTime"`
	PackageVersion                *string            `json:"packageVersion,omitempty"`
	PackageAuthor                 *string            `json:"packageAuthor,omitempty"`
	PackageUpdateMechanism        UpdateMechanism    `json:"packageUpdateMechanism,omitempty"`
	PackageUpdateMechanismMessage *string            `json:"packageUpdateMechanismMessage,omitempty"`
}

// Task is a scheduled background job in Radarr.
type Task struct {
	ID            int       `json:"id"`
	Name          *string   `json:"name,omitempty"`
	TaskName      *string   `json:"taskName,omitempty"`
	Interval      int       `json:"interval"`
	LastExecution time.Time `json:"lastExecution"`
	LastStartTime time.Time `json:"lastStartTime"`
	NextExecution time.Time `json:"nextExecution"`
	LastDuration  *string   `json:"lastDuration,omitempty"`
}

// BackupType classifies what triggered a backup.
type BackupType string

// Backup type values.
const (
	BackupTypeScheduled BackupType = "scheduled"
	BackupTypeManual    BackupType = "manual"
	BackupTypeUpdate    BackupType = "update"
)

// Backup represents a Radarr database backup file.
type Backup struct {
	ID   int        `json:"id"`
	Name *string    `json:"name,omitempty"`
	Path *string    `json:"path,omitempty"`
	Type BackupType `json:"type,omitempty"`
	Size int64      `json:"size"`
	Time time.Time  `json:"time"`
}

// HealthCheckResult is the severity of a health check finding.
type HealthCheckResult string

// Health check result values.
const (
	HealthCheckResultOk      HealthCheckResult = "ok"
	HealthCheckResultNotice  HealthCheckResult = "notice"
	HealthCheckResultWarning HealthCheckResult = "warning"
	HealthCheckResultError   HealthCheckResult = "error"
)

// HealthRecord is a single health check result reported by Radarr.
type HealthRecord struct {
	ID      int               `json:"id"`
	Source  *string           `json:"source,omitempty"`
	Type    HealthCheckResult `json:"type,omitempty"`
	Message *string           `json:"message,omitempty"`
	WikiURL *string           `json:"wikiUrl,omitempty"`
}

// DiskSpace reports free and total space for a configured media path.
type DiskSpace struct {
	ID         int     `json:"id"`
	Path       *string `json:"path,omitempty"`
	Label      *string `json:"label,omitempty"`
	FreeSpace  int64   `json:"freeSpace"`
	TotalSpace int64   `json:"totalSpace"`
}

// UpdateChanges categorises new and fixed items in a Radarr release.
type UpdateChanges struct {
	New   []string `json:"new,omitempty"`
	Fixed []string `json:"fixed,omitempty"`
}

// Update represents an available or installed Radarr application update.
type Update struct {
	ID          int            `json:"id"`
	Version     *string        `json:"version,omitempty"`
	Branch      *string        `json:"branch,omitempty"`
	ReleaseDate time.Time      `json:"releaseDate"`
	FileName    *string        `json:"fileName,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Installed   bool           `json:"installed"`
	InstalledOn *time.Time     `json:"installedOn,omitempty"`
	Installable bool           `json:"installable"`
	Latest      bool           `json:"latest"`
	Changes     *UpdateChanges `json:"changes,omitempty"`
	Hash        *string        `json:"hash,omitempty"`
}

// ---------------------------------------------------------------------------
// Status
// ---------------------------------------------------------------------------

// Status returns runtime and environment information about the Radarr instance.
func (s *SystemService) Status(ctx context.Context) (*SystemStatus, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&SystemStatus{}).
		Get("/api/v3/system/status")
	if err != nil {
		return nil, fmt.Errorf("radarr: get system status: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get system status: %w", err)
	}

	return ptrResult[SystemStatus](resp)
}

// ---------------------------------------------------------------------------
// Tasks
// ---------------------------------------------------------------------------

// ListTasks returns all scheduled background tasks.
func (s *SystemService) ListTasks(ctx context.Context) ([]Task, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Task{}).
		Get("/api/v3/system/task")
	if err != nil {
		return nil, fmt.Errorf("radarr: list system tasks: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list system tasks: %w", err)
	}

	return derefResult[[]Task](resp)
}

// GetTask returns a single scheduled task by ID.
func (s *SystemService) GetTask(ctx context.Context, id int) (*Task, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&Task{}).
		SetPathParam("id", strconv.Itoa(id)).
		Get("/api/v3/system/task/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: get system task %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get system task %d: %w", id, err)
	}

	return ptrResult[Task](resp)
}

// ---------------------------------------------------------------------------
// Backup
// ---------------------------------------------------------------------------

// ListBackups returns all available database backup files.
func (s *SystemService) ListBackups(ctx context.Context) ([]Backup, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Backup{}).
		Get("/api/v3/system/backup")
	if err != nil {
		return nil, fmt.Errorf("radarr: list backups: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: list backups: %w", err)
	}

	return derefResult[[]Backup](resp)
}

// DeleteBackup removes a backup file by ID.
func (s *SystemService) DeleteBackup(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Delete("/api/v3/system/backup/{id}")
	if err != nil {
		return fmt.Errorf("radarr: delete backup %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: delete backup %d: %w", id, err)
	}

	return nil
}

// RestoreBackup triggers a restore from a backup file by ID.
func (s *SystemService) RestoreBackup(ctx context.Context, id int) error {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("id", strconv.Itoa(id)).
		Post("/api/v3/system/backup/restore/{id}")
	if err != nil {
		return fmt.Errorf("radarr: restore backup %d: %w", id, err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: restore backup %d: %w", id, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Restart / Shutdown
// ---------------------------------------------------------------------------

// Restart sends a restart command to the Radarr instance.
func (s *SystemService) Restart(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/system/restart")
	if err != nil {
		return fmt.Errorf("radarr: restart: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: restart: %w", err)
	}

	return nil
}

// Shutdown sends a shutdown command to the Radarr instance.
func (s *SystemService) Shutdown(ctx context.Context) error {
	resp, err := s.client.R().
		SetContext(ctx).
		Post("/api/v3/system/shutdown")
	if err != nil {
		return fmt.Errorf("radarr: shutdown: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return fmt.Errorf("radarr: shutdown: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

// Health returns the current health check results for the Radarr instance.
func (s *SystemService) Health(ctx context.Context) ([]HealthRecord, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]HealthRecord{}).
		Get("/api/v3/health")
	if err != nil {
		return nil, fmt.Errorf("radarr: get health: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get health: %w", err)
	}

	return derefResult[[]HealthRecord](resp)
}

// ---------------------------------------------------------------------------
// DiskSpace
// ---------------------------------------------------------------------------

// DiskSpace returns disk usage information for configured media paths.
func (s *SystemService) DiskSpace(ctx context.Context) ([]DiskSpace, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]DiskSpace{}).
		Get("/api/v3/diskspace")
	if err != nil {
		return nil, fmt.Errorf("radarr: get disk space: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get disk space: %w", err)
	}

	return derefResult[[]DiskSpace](resp)
}

// ---------------------------------------------------------------------------
// Updates
// ---------------------------------------------------------------------------

// Updates returns available application updates for the Radarr instance.
func (s *SystemService) Updates(ctx context.Context) ([]Update, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&[]Update{}).
		Get("/api/v3/update")
	if err != nil {
		return nil, fmt.Errorf("radarr: get updates: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get updates: %w", err)
	}

	return derefResult[[]Update](resp)
}
