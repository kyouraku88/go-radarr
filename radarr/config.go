package radarr

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// ConfigService provides methods for all /config/* endpoints.
type ConfigService service

// ---------------------------------------------------------------------------
// Host config
// ---------------------------------------------------------------------------

// HostConfig holds the Radarr host and network configuration.
type HostConfig struct {
	ID                        int                        `json:"id"`
	BindAddress               *string                    `json:"bindAddress,omitempty"`
	Port                      int                        `json:"port"`
	SslPort                   int                        `json:"sslPort"`
	EnableSsl                 bool                       `json:"enableSsl"`
	LaunchBrowser             bool                       `json:"launchBrowser"`
	AuthenticationMethod      AuthenticationType         `json:"authenticationMethod,omitempty"`
	AuthenticationRequired    AuthenticationRequiredType `json:"authenticationRequired,omitempty"`
	AnalyticsEnabled          bool                       `json:"analyticsEnabled"`
	Username                  *string                    `json:"username,omitempty"`
	Password                  *string                    `json:"password,omitempty"`
	PasswordConfirmation      *string                    `json:"passwordConfirmation,omitempty"`
	LogLevel                  *string                    `json:"logLevel,omitempty"`
	LogSizeLimit              int                        `json:"logSizeLimit"`
	ConsoleLogLevel           *string                    `json:"consoleLogLevel,omitempty"`
	Branch                    *string                    `json:"branch,omitempty"`
	APIKey                    *string                    `json:"apiKey,omitempty"`
	SslCertPath               *string                    `json:"sslCertPath,omitempty"`
	SslCertPassword           *string                    `json:"sslCertPassword,omitempty"`
	URLBase                   *string                    `json:"urlBase,omitempty"`
	InstanceName              *string                    `json:"instanceName,omitempty"`
	ApplicationURL            *string                    `json:"applicationUrl,omitempty"`
	UpdateAutomatically       bool                       `json:"updateAutomatically"`
	UpdateMechanism           UpdateMechanism            `json:"updateMechanism,omitempty"`
	UpdateScriptPath          *string                    `json:"updateScriptPath,omitempty"`
	ProxyEnabled              bool                       `json:"proxyEnabled"`
	ProxyType                 ProxyType                  `json:"proxyType,omitempty"`
	ProxyHostname             *string                    `json:"proxyHostname,omitempty"`
	ProxyPort                 int                        `json:"proxyPort"`
	ProxyUsername             *string                    `json:"proxyUsername,omitempty"`
	ProxyPassword             *string                    `json:"proxyPassword,omitempty"`
	ProxyBypassFilter         *string                    `json:"proxyBypassFilter,omitempty"`
	ProxyBypassLocalAddresses bool                       `json:"proxyBypassLocalAddresses"`
	CertificateValidation     CertificateValidationType  `json:"certificateValidation,omitempty"`
	BackupFolder              *string                    `json:"backupFolder,omitempty"`
	BackupInterval            int                        `json:"backupInterval"`
	BackupRetention           int                        `json:"backupRetention"`
	TrustCgnatIPAddresses     bool                       `json:"trustCgnatIpAddresses"`
}

// GetHostConfig returns the host configuration.
func (s *ConfigService) GetHostConfig(ctx context.Context) (*HostConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&HostConfig{}).
		Get("/api/v3/config/host")
	if err != nil {
		return nil, fmt.Errorf("radarr: get host config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get host config: %w", err)
	}

	return ptrResult[HostConfig](resp)
}

// UpdateHostConfig replaces the host configuration by ID.
func (s *ConfigService) UpdateHostConfig(ctx context.Context, id int, body HostConfig) (*HostConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&HostConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/host/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update host config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update host config: %w", err)
	}

	return ptrResult[HostConfig](resp)
}

// ---------------------------------------------------------------------------
// UI config
// ---------------------------------------------------------------------------

// UIConfig holds the Radarr user interface configuration.
type UIConfig struct {
	ID                       int                    `json:"id"`
	FirstDayOfWeek           int                    `json:"firstDayOfWeek"`
	CalendarWeekColumnHeader *string                `json:"calendarWeekColumnHeader,omitempty"`
	MovieRuntimeFormat       MovieRuntimeFormatType `json:"movieRuntimeFormat,omitempty"`
	ShortDateFormat          *string                `json:"shortDateFormat,omitempty"`
	LongDateFormat           *string                `json:"longDateFormat,omitempty"`
	TimeFormat               *string                `json:"timeFormat,omitempty"`
	ShowRelativeDates        bool                   `json:"showRelativeDates"`
	EnableColorImpairedMode  bool                   `json:"enableColorImpairedMode"`
	MovieInfoLanguage        int                    `json:"movieInfoLanguage"`
	UILanguage               int                    `json:"uiLanguage"`
	Theme                    *string                `json:"theme,omitempty"`
}

// GetUIConfig returns the UI configuration.
func (s *ConfigService) GetUIConfig(ctx context.Context) (*UIConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&UIConfig{}).
		Get("/api/v3/config/ui")
	if err != nil {
		return nil, fmt.Errorf("radarr: get ui config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get ui config: %w", err)
	}

	return ptrResult[UIConfig](resp)
}

// UpdateUIConfig replaces the UI configuration by ID.
func (s *ConfigService) UpdateUIConfig(ctx context.Context, id int, body UIConfig) (*UIConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&UIConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/ui/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update ui config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update ui config: %w", err)
	}

	return ptrResult[UIConfig](resp)
}

// ---------------------------------------------------------------------------
// Naming config
// ---------------------------------------------------------------------------

// NamingConfig holds the file and folder naming configuration.
type NamingConfig struct {
	ID                       int                    `json:"id"`
	RenameMovies             bool                   `json:"renameMovies"`
	ReplaceIllegalCharacters bool                   `json:"replaceIllegalCharacters"`
	ColonReplacementFormat   ColonReplacementFormat `json:"colonReplacementFormat,omitempty"`
	StandardMovieFormat      *string                `json:"standardMovieFormat,omitempty"`
	MovieFolderFormat        *string                `json:"movieFolderFormat,omitempty"`
}

// GetNamingConfig returns the naming configuration.
func (s *ConfigService) GetNamingConfig(ctx context.Context) (*NamingConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&NamingConfig{}).
		Get("/api/v3/config/naming")
	if err != nil {
		return nil, fmt.Errorf("radarr: get naming config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get naming config: %w", err)
	}

	return ptrResult[NamingConfig](resp)
}

// UpdateNamingConfig replaces the naming configuration by ID.
func (s *ConfigService) UpdateNamingConfig(ctx context.Context, id int, body NamingConfig) (*NamingConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&NamingConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/naming/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update naming config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update naming config: %w", err)
	}

	return ptrResult[NamingConfig](resp)
}

// NamingExamplesOption is a functional option for ConfigService.GetNamingExamples.
type NamingExamplesOption func(*resty.Request)

// WithNamingExamplesRenameMovies sets the renameMovies flag for example generation.
func WithNamingExamplesRenameMovies(v bool) NamingExamplesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("renameMovies", strconv.FormatBool(v))
	}
}

// WithNamingExamplesReplaceIllegalCharacters sets the replaceIllegalCharacters flag.
func WithNamingExamplesReplaceIllegalCharacters(v bool) NamingExamplesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("replaceIllegalCharacters", strconv.FormatBool(v))
	}
}

// WithNamingExamplesStandardMovieFormat sets the standard movie format for example generation.
func WithNamingExamplesStandardMovieFormat(format string) NamingExamplesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("standardMovieFormat", format)
	}
}

// WithNamingExamplesMovieFolderFormat sets the movie folder format for example generation.
func WithNamingExamplesMovieFolderFormat(format string) NamingExamplesOption {
	return func(r *resty.Request) {
		r.SetQueryParam("movieFolderFormat", format)
	}
}

// GetNamingExamples returns example filenames generated by the given naming config.
func (s *ConfigService) GetNamingExamples(ctx context.Context, opts ...NamingExamplesOption) (map[string]string, error) {
	req := s.client.R().
		SetContext(ctx)
	for _, o := range opts {
		o(req)
	}

	resp, err := req.Get("/api/v3/config/naming/examples")
	if err != nil {
		return nil, fmt.Errorf("radarr: get naming examples: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get naming examples: %w", err)
	}

	var result map[string]string
	if jsonErr := json.Unmarshal(resp.Body(), &result); jsonErr != nil {
		return nil, fmt.Errorf("radarr: get naming examples: %w", jsonErr)
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Media management config
// ---------------------------------------------------------------------------

// MediaManagementConfig holds the media file management configuration.
type MediaManagementConfig struct {
	ID                                      int                    `json:"id"`
	AutoUnmonitorPreviouslyDownloadedMovies bool                   `json:"autoUnmonitorPreviouslyDownloadedMovies"`
	RecycleBin                              *string                `json:"recycleBin,omitempty"`
	RecycleBinCleanupDays                   int                    `json:"recycleBinCleanupDays"`
	DownloadPropersAndRepacks               ProperDownloadTypes    `json:"downloadPropersAndRepacks,omitempty"`
	CreateEmptyMovieFolders                 bool                   `json:"createEmptyMovieFolders"`
	DeleteEmptyFolders                      bool                   `json:"deleteEmptyFolders"`
	FileDate                                FileDateType           `json:"fileDate,omitempty"`
	RescanAfterRefresh                      RescanAfterRefreshType `json:"rescanAfterRefresh,omitempty"`
	AutoRenameFolders                       bool                   `json:"autoRenameFolders"`
	PathsDefaultStatic                      bool                   `json:"pathsDefaultStatic"`
	SetPermissionsLinux                     bool                   `json:"setPermissionsLinux"`
	ChmodFolder                             *string                `json:"chmodFolder,omitempty"`
	ChownGroup                              *string                `json:"chownGroup,omitempty"`
	SkipFreeSpaceCheckWhenImporting         bool                   `json:"skipFreeSpaceCheckWhenImporting"`
	MinimumFreeSpaceWhenImporting           int                    `json:"minimumFreeSpaceWhenImporting"`
	CopyUsingHardlinks                      bool                   `json:"copyUsingHardlinks"`
	UseScriptImport                         bool                   `json:"useScriptImport"`
	ScriptImportPath                        *string                `json:"scriptImportPath,omitempty"`
	ImportExtraFiles                        bool                   `json:"importExtraFiles"`
	ExtraFileExtensions                     *string                `json:"extraFileExtensions,omitempty"`
	EnableMediaInfo                         bool                   `json:"enableMediaInfo"`
}

// GetMediaManagementConfig returns the media management configuration.
func (s *ConfigService) GetMediaManagementConfig(ctx context.Context) (*MediaManagementConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MediaManagementConfig{}).
		Get("/api/v3/config/mediamanagement")
	if err != nil {
		return nil, fmt.Errorf("radarr: get media management config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get media management config: %w", err)
	}

	return ptrResult[MediaManagementConfig](resp)
}

// UpdateMediaManagementConfig replaces the media management configuration by ID.
func (s *ConfigService) UpdateMediaManagementConfig(ctx context.Context, id int, body MediaManagementConfig) (*MediaManagementConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MediaManagementConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/mediamanagement/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update media management config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update media management config: %w", err)
	}

	return ptrResult[MediaManagementConfig](resp)
}

// ---------------------------------------------------------------------------
// Download client config
// ---------------------------------------------------------------------------

// DownloadClientConfig holds the download client integration configuration.
type DownloadClientConfig struct {
	ID                                        int     `json:"id"`
	DownloadClientWorkingFolders              *string `json:"downloadClientWorkingFolders,omitempty"`
	EnableCompletedDownloadHandling           bool    `json:"enableCompletedDownloadHandling"`
	CheckForFinishedDownloadInterval          int     `json:"checkForFinishedDownloadInterval"`
	AutoRedownloadFailed                      bool    `json:"autoRedownloadFailed"`
	AutoRedownloadFailedFromInteractiveSearch bool    `json:"autoRedownloadFailedFromInteractiveSearch"`
}

// GetDownloadClientConfig returns the download client configuration.
func (s *ConfigService) GetDownloadClientConfig(ctx context.Context) (*DownloadClientConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DownloadClientConfig{}).
		Get("/api/v3/config/downloadclient")
	if err != nil {
		return nil, fmt.Errorf("radarr: get download client config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get download client config: %w", err)
	}

	return ptrResult[DownloadClientConfig](resp)
}

// UpdateDownloadClientConfig replaces the download client configuration by ID.
func (s *ConfigService) UpdateDownloadClientConfig(ctx context.Context, id int, body DownloadClientConfig) (*DownloadClientConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&DownloadClientConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/downloadclient/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update download client config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update download client config: %w", err)
	}

	return ptrResult[DownloadClientConfig](resp)
}

// ---------------------------------------------------------------------------
// Import list config
// ---------------------------------------------------------------------------

// ImportListConfig holds the import list synchronisation configuration.
type ImportListConfig struct {
	ID            int     `json:"id"`
	ListSyncLevel *string `json:"listSyncLevel,omitempty"`
}

// GetImportListConfig returns the import list configuration.
func (s *ConfigService) GetImportListConfig(ctx context.Context) (*ImportListConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportListConfig{}).
		Get("/api/v3/config/importlist")
	if err != nil {
		return nil, fmt.Errorf("radarr: get import list config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get import list config: %w", err)
	}

	return ptrResult[ImportListConfig](resp)
}

// UpdateImportListConfig replaces the import list configuration by ID.
func (s *ConfigService) UpdateImportListConfig(ctx context.Context, id int, body ImportListConfig) (*ImportListConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&ImportListConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/importlist/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update import list config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update import list config: %w", err)
	}

	return ptrResult[ImportListConfig](resp)
}

// ---------------------------------------------------------------------------
// Indexer config
// ---------------------------------------------------------------------------

// IndexerConfig holds the indexer global configuration.
type IndexerConfig struct {
	ID                       int     `json:"id"`
	MinimumAge               int     `json:"minimumAge"`
	MaximumSize              int     `json:"maximumSize"`
	Retention                int     `json:"retention"`
	RssSyncInterval          int     `json:"rssSyncInterval"`
	PreferIndexerFlags       bool    `json:"preferIndexerFlags"`
	AvailabilityDelay        int     `json:"availabilityDelay"`
	AllowHardcodedSubs       bool    `json:"allowHardcodedSubs"`
	WhitelistedHardcodedSubs *string `json:"whitelistedHardcodedSubs,omitempty"`
}

// GetIndexerConfig returns the indexer configuration.
func (s *ConfigService) GetIndexerConfig(ctx context.Context) (*IndexerConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&IndexerConfig{}).
		Get("/api/v3/config/indexer")
	if err != nil {
		return nil, fmt.Errorf("radarr: get indexer config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get indexer config: %w", err)
	}

	return ptrResult[IndexerConfig](resp)
}

// UpdateIndexerConfig replaces the indexer configuration by ID.
func (s *ConfigService) UpdateIndexerConfig(ctx context.Context, id int, body IndexerConfig) (*IndexerConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&IndexerConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/indexer/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update indexer config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update indexer config: %w", err)
	}

	return ptrResult[IndexerConfig](resp)
}

// ---------------------------------------------------------------------------
// Metadata config
// ---------------------------------------------------------------------------

// MetadataConfig holds the metadata global configuration.
type MetadataConfig struct {
	ID                   int             `json:"id"`
	CertificationCountry TMDbCountryCode `json:"certificationCountry,omitempty"`
}

// GetMetadataConfig returns the metadata configuration.
func (s *ConfigService) GetMetadataConfig(ctx context.Context) (*MetadataConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MetadataConfig{}).
		Get("/api/v3/config/metadata")
	if err != nil {
		return nil, fmt.Errorf("radarr: get metadata config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: get metadata config: %w", err)
	}

	return ptrResult[MetadataConfig](resp)
}

// UpdateMetadataConfig replaces the metadata configuration by ID.
func (s *ConfigService) UpdateMetadataConfig(ctx context.Context, id int, body MetadataConfig) (*MetadataConfig, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&MetadataConfig{}).
		SetPathParam("id", strconv.Itoa(id)).
		SetBody(body).
		Put("/api/v3/config/metadata/{id}")
	if err != nil {
		return nil, fmt.Errorf("radarr: update metadata config: %w", err)
	}

	if err := checkResponse(resp); err != nil {
		return nil, fmt.Errorf("radarr: update metadata config: %w", err)
	}

	return ptrResult[MetadataConfig](resp)
}
