// Package radarr provides a Go client for the Radarr v3 API.
package radarr

import (
	"context"
	"errors"
	"iter"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// Radarr is the top-level client. Use New to construct one.
type Radarr struct {
	common service

	AltTitle          *AltTitleService
	AutoTagging       *AutoTaggingService
	Blocklist         *BlocklistService
	Calendar          *CalendarService
	Collection        *CollectionService
	Command           *CommandService
	Config            *ConfigService
	Credit            *CreditService
	CustomFilter      *CustomFilterService
	CustomFormat      *CustomFormatService
	DelayProfile      *DelayProfileService
	DownloadClient    *DownloadClientService
	Exclusions        *ExclusionsService
	ExtraFile         *ExtraFileService
	Filesystem        *FilesystemService
	History           *HistoryService
	ImportList        *ImportListService
	Indexer           *IndexerService
	IndexerFlag       *IndexerFlagService
	Language          *LanguageService
	Localization      *LocalizationService
	Log               *LogService
	ManualImport      *ManualImportService
	Metadata          *MetadataService
	Movie             *MovieService
	MovieEditor       *MovieEditorService
	MovieFile         *MovieFileService
	Notification      *NotificationService
	Parse             *ParseService
	QualityDefinition *QualityDefinitionService
	QualityProfile    *QualityProfileService
	Queue             *QueueService
	Release           *ReleaseService
	ReleaseProfile    *ReleaseProfileService
	RemotePathMapping *RemotePathMappingService
	Rename            *RenameService
	RootFolder        *RootFolderService
	System            *SystemService
	Tag               *TagService
	Wanted            *WantedService
}

type service struct {
	client *resty.Client
}

// ClientOption configures the client during construction.
type ClientOption func(*clientConfig) error

type clientConfig struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// WithBaseURL sets the base URL of the Radarr instance (e.g. "http://localhost:7878").
func WithBaseURL(url string) ClientOption {
	return func(c *clientConfig) error {
		c.baseURL = url
		return nil
	}
}

// WithAPIKey sets the API key used for authentication via the X-Api-Key header.
func WithAPIKey(key string) ClientOption {
	return func(c *clientConfig) error {
		c.apiKey = key
		return nil
	}
}

// WithHTTPClient replaces the default HTTP client (e.g. to set timeouts or a custom transport).
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *clientConfig) error {
		c.httpClient = hc
		return nil
	}
}

// New creates a new Radarr client. WithBaseURL and WithAPIKey are required.
func New(opts ...ClientOption) (*Radarr, error) {
	cfg := &clientConfig{}
	for _, o := range opts {
		if err := o(cfg); err != nil {
			return nil, err
		}
	}

	if cfg.baseURL == "" {
		return nil, errors.New("radarr: base URL is required")
	}

	if cfg.apiKey == "" {
		return nil, errors.New("radarr: API key is required")
	}

	rc := resty.New().
		SetBaseURL(cfg.baseURL).
		SetHeader("X-Api-Key", cfg.apiKey)
	if cfg.httpClient != nil {
		rc.SetTransport(cfg.httpClient.Transport)
	}

	r := &Radarr{}
	r.common.client = rc
	r.AltTitle = (*AltTitleService)(&r.common)
	r.AutoTagging = (*AutoTaggingService)(&r.common)
	r.Blocklist = (*BlocklistService)(&r.common)
	r.Calendar = (*CalendarService)(&r.common)
	r.Collection = (*CollectionService)(&r.common)
	r.Command = (*CommandService)(&r.common)
	r.Config = (*ConfigService)(&r.common)
	r.Credit = (*CreditService)(&r.common)
	r.CustomFilter = (*CustomFilterService)(&r.common)
	r.CustomFormat = (*CustomFormatService)(&r.common)
	r.DelayProfile = (*DelayProfileService)(&r.common)
	r.DownloadClient = (*DownloadClientService)(&r.common)
	r.Exclusions = (*ExclusionsService)(&r.common)
	r.ExtraFile = (*ExtraFileService)(&r.common)
	r.Filesystem = (*FilesystemService)(&r.common)
	r.History = (*HistoryService)(&r.common)
	r.ImportList = (*ImportListService)(&r.common)
	r.Indexer = (*IndexerService)(&r.common)
	r.IndexerFlag = (*IndexerFlagService)(&r.common)
	r.Language = (*LanguageService)(&r.common)
	r.Localization = (*LocalizationService)(&r.common)
	r.Log = (*LogService)(&r.common)
	r.ManualImport = (*ManualImportService)(&r.common)
	r.Metadata = (*MetadataService)(&r.common)
	r.Movie = (*MovieService)(&r.common)
	r.MovieEditor = (*MovieEditorService)(&r.common)
	r.MovieFile = (*MovieFileService)(&r.common)
	r.Notification = (*NotificationService)(&r.common)
	r.Parse = (*ParseService)(&r.common)
	r.QualityDefinition = (*QualityDefinitionService)(&r.common)
	r.QualityProfile = (*QualityProfileService)(&r.common)
	r.Queue = (*QueueService)(&r.common)
	r.Release = (*ReleaseService)(&r.common)
	r.ReleaseProfile = (*ReleaseProfileService)(&r.common)
	r.RemotePathMapping = (*RemotePathMappingService)(&r.common)
	r.Rename = (*RenameService)(&r.common)
	r.RootFolder = (*RootFolderService)(&r.common)
	r.System = (*SystemService)(&r.common)
	r.Tag = (*TagService)(&r.common)
	r.Wanted = (*WantedService)(&r.common)

	return r, nil
}

// ---------------------------------------------------------------------------
// Shared enums
// ---------------------------------------------------------------------------

// SortDirection controls the order of paginated results.
type SortDirection string

// Sort direction values.
const (
	SortDirectionDefault    SortDirection = "default"
	SortDirectionAscending  SortDirection = "ascending"
	SortDirectionDescending SortDirection = "descending"
)

// DownloadProtocol identifies the protocol used by a download client.
type DownloadProtocol string

// Download protocol values.
const (
	DownloadProtocolUnknown DownloadProtocol = "unknown"
	DownloadProtocolUsenet  DownloadProtocol = "usenet"
	DownloadProtocolTorrent DownloadProtocol = "torrent"
)

// MediaCoverType identifies the kind of artwork image.
type MediaCoverType string

// Media cover type values.
const (
	MediaCoverTypeUnknown    MediaCoverType = "unknown"
	MediaCoverTypePoster     MediaCoverType = "poster"
	MediaCoverTypeBanner     MediaCoverType = "banner"
	MediaCoverTypeFanart     MediaCoverType = "fanart"
	MediaCoverTypeScreenshot MediaCoverType = "screenshot"
	MediaCoverTypeHeadshot   MediaCoverType = "headshot"
	MediaCoverTypeClearLogo  MediaCoverType = "clearlogo"
)

// MovieStatusType represents a movie's release lifecycle state.
type MovieStatusType string

// Movie status type values.
const (
	MovieStatusTypeTBA       MovieStatusType = "tba"
	MovieStatusTypeAnnounced MovieStatusType = "announced"
	MovieStatusTypeInCinemas MovieStatusType = "inCinemas"
	MovieStatusTypeReleased  MovieStatusType = "released"
	MovieStatusTypeDeleted   MovieStatusType = "deleted"
)

// MonitorTypes controls what content Radarr monitors for a movie.
type MonitorTypes string

// Monitor type values.
const (
	MonitorTypesMovieOnly          MonitorTypes = "movieOnly"
	MonitorTypesMovieAndCollection MonitorTypes = "movieAndCollection"
	MonitorTypesNone               MonitorTypes = "none"
)

// AddMovieMethod describes how a movie was added to the library.
type AddMovieMethod string

// Add movie method values.
const (
	AddMovieMethodManual     AddMovieMethod = "manual"
	AddMovieMethodList       AddMovieMethod = "list"
	AddMovieMethodCollection AddMovieMethod = "collection"
)

// RatingType distinguishes audience from critic ratings.
type RatingType string

// Rating type values.
const (
	RatingTypeUser   RatingType = "user"
	RatingTypeCritic RatingType = "critic"
)

// QualitySource identifies the release medium of a quality tier.
type QualitySource string

// Quality source values.
const (
	QualitySourceUnknown   QualitySource = "unknown"
	QualitySourceCam       QualitySource = "cam"
	QualitySourceTelesync  QualitySource = "telesync"
	QualitySourceTelecine  QualitySource = "telecine"
	QualitySourceWorkprint QualitySource = "workprint"
	QualitySourceDVD       QualitySource = "dvd"
	QualitySourceTV        QualitySource = "tv"
	QualitySourceWebDL     QualitySource = "webdl"
	QualitySourceWebRip    QualitySource = "webrip"
	QualitySourceBluray    QualitySource = "bluray"
)

// Modifier further qualifies a quality source.
type Modifier string

// Modifier values.
const (
	ModifierNone     Modifier = "none"
	ModifierRegional Modifier = "regional"
	ModifierScreener Modifier = "screener"
	ModifierRawHD    Modifier = "rawhd"
	ModifierBRDisk   Modifier = "brdisk"
	ModifierRemux    Modifier = "remux"
)

// SourceType identifies where a movie's metadata originated.
type SourceType string

// Source type values.
const (
	SourceTypeTmdb     SourceType = "tmdb"
	SourceTypeMappings SourceType = "mappings"
	SourceTypeUser     SourceType = "user"
	SourceTypeIndexer  SourceType = "indexer"
)

// AuthenticationRequiredType specifies when authentication is required.
type AuthenticationRequiredType string

// Authentication required type values.
const (
	AuthenticationRequiredTypeEnabled                   AuthenticationRequiredType = "enabled"
	AuthenticationRequiredTypeDisabledForLocalAddresses AuthenticationRequiredType = "disabledForLocalAddresses"
)

// ProxyType identifies the proxy protocol.
type ProxyType string

// Proxy type values.
const (
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeSocks4 ProxyType = "socks4"
	ProxyTypeSocks5 ProxyType = "socks5"
)

// CertificateValidationType controls TLS certificate validation behaviour.
type CertificateValidationType string

// Certificate validation type values.
const (
	CertificateValidationTypeEnabled                   CertificateValidationType = "enabled"
	CertificateValidationTypeDisabledForLocalAddresses CertificateValidationType = "disabledForLocalAddresses"
	CertificateValidationTypeDisabled                  CertificateValidationType = "disabled"
)

// MovieRuntimeFormatType controls how runtimes are displayed.
type MovieRuntimeFormatType string

// Movie runtime format type values.
const (
	MovieRuntimeFormatTypeHoursMinutes MovieRuntimeFormatType = "hoursMinutes"
	MovieRuntimeFormatTypeMinutes      MovieRuntimeFormatType = "minutes"
)

// ColonReplacementFormat controls how colons in filenames are handled.
type ColonReplacementFormat string

// Colon replacement format values.
const (
	ColonReplacementFormatDelete         ColonReplacementFormat = "delete"
	ColonReplacementFormatDash           ColonReplacementFormat = "dash"
	ColonReplacementFormatSpaceDash      ColonReplacementFormat = "spaceDash"
	ColonReplacementFormatSpaceDashSpace ColonReplacementFormat = "spaceDashSpace"
	ColonReplacementFormatSmart          ColonReplacementFormat = "smart"
)

// ProperDownloadTypes controls how Radarr handles proper/repack releases.
type ProperDownloadTypes string

// Proper download type values.
const (
	ProperDownloadTypesPreferAndUpgrade ProperDownloadTypes = "preferAndUpgrade"
	ProperDownloadTypesDoNotUpgrade     ProperDownloadTypes = "doNotUpgrade"
	ProperDownloadTypesDoNotPrefer      ProperDownloadTypes = "doNotPrefer"
)

// FileDateType controls which date is used as the file modification date.
type FileDateType string

// File date type values.
const (
	FileDateTypeNone    FileDateType = "none"
	FileDateTypeCinemas FileDateType = "cinemas"
	FileDateTypeRelease FileDateType = "release"
)

// RescanAfterRefreshType controls when Radarr rescans the disk after a refresh.
type RescanAfterRefreshType string

// Rescan after refresh type values.
const (
	RescanAfterRefreshTypeAlways      RescanAfterRefreshType = "always"
	RescanAfterRefreshTypeAfterManual RescanAfterRefreshType = "afterManual"
	RescanAfterRefreshTypeNever       RescanAfterRefreshType = "never"
)

// TMDbCountryCode selects the content rating country for TMDb metadata.
type TMDbCountryCode string

// TMDb country code values.
const (
	TMDbCountryCodeAU TMDbCountryCode = "au"
	TMDbCountryCodeBR TMDbCountryCode = "br"
	TMDbCountryCodeCA TMDbCountryCode = "ca"
	TMDbCountryCodeFR TMDbCountryCode = "fr"
	TMDbCountryCodeDE TMDbCountryCode = "de"
	TMDbCountryCodeGB TMDbCountryCode = "gb"
	TMDbCountryCodeIN TMDbCountryCode = "in"
	TMDbCountryCodeIE TMDbCountryCode = "ie"
	TMDbCountryCodeIT TMDbCountryCode = "it"
	TMDbCountryCodeNZ TMDbCountryCode = "nz"
	TMDbCountryCodeRO TMDbCountryCode = "ro"
	TMDbCountryCodeES TMDbCountryCode = "es"
	TMDbCountryCodeUS TMDbCountryCode = "us"
)

// ExtraFileType identifies the kind of extra file alongside a movie.
type ExtraFileType string

// Extra file type values.
const (
	ExtraFileTypeSubtitle ExtraFileType = "subtitle"
	ExtraFileTypeMetadata ExtraFileType = "metadata"
	ExtraFileTypeOther    ExtraFileType = "other"
)

// ImportListType categorises the source of an import list.
type ImportListType string

// Import list type values.
const (
	ImportListTypeProgram  ImportListType = "program"
	ImportListTypeTMDb     ImportListType = "tmdb"
	ImportListTypeTrakt    ImportListType = "trakt"
	ImportListTypePlex     ImportListType = "plex"
	ImportListTypeSimkl    ImportListType = "simkl"
	ImportListTypeOther    ImportListType = "other"
	ImportListTypeAdvanced ImportListType = "advanced"
)

// RejectionType indicates whether a rejection is permanent or temporary.
type RejectionType string

// Rejection type values.
const (
	RejectionTypePermanent RejectionType = "permanent"
	RejectionTypeTemporary RejectionType = "temporary"
)

// PrivacyLevel controls how sensitive a provider field is.
type PrivacyLevel string

// Privacy level values.
const (
	PrivacyLevelNormal   PrivacyLevel = "normal"
	PrivacyLevelPassword PrivacyLevel = "password"
	PrivacyLevelAPIKey   PrivacyLevel = "apiKey"
	PrivacyLevelUserName PrivacyLevel = "userName"
)

// ProviderMessageType is the severity of a provider status message.
type ProviderMessageType string

// Provider message type values.
const (
	ProviderMessageTypeInfo    ProviderMessageType = "info"
	ProviderMessageTypeWarning ProviderMessageType = "warning"
	ProviderMessageTypeError   ProviderMessageType = "error"
)

// ApplyTags controls how tags are merged during bulk updates.
type ApplyTags string

// Apply tags values.
const (
	ApplyTagsAdd     ApplyTags = "add"
	ApplyTagsRemove  ApplyTags = "remove"
	ApplyTagsReplace ApplyTags = "replace"
)

// CreditType distinguishes cast members from crew members.
type CreditType string

// Credit type values.
const (
	CreditTypeCast CreditType = "cast"
	CreditTypeCrew CreditType = "crew"
)

// ---------------------------------------------------------------------------
// Shared value types
// ---------------------------------------------------------------------------

// SelectOption is a selectable item in a provider configuration dropdown.
type SelectOption struct {
	Value        int     `json:"value"`
	Name         *string `json:"name,omitempty"`
	Order        int     `json:"order"`
	Hint         *string `json:"hint,omitempty"`
	DividerAfter bool    `json:"dividerAfter"`
}

// Field is a single configuration field within a provider resource.
type Field struct {
	Order                       int            `json:"order"`
	Name                        *string        `json:"name,omitempty"`
	Label                       *string        `json:"label,omitempty"`
	Unit                        *string        `json:"unit,omitempty"`
	HelpText                    *string        `json:"helpText,omitempty"`
	HelpTextWarning             *string        `json:"helpTextWarning,omitempty"`
	HelpLink                    *string        `json:"helpLink,omitempty"`
	Value                       any            `json:"value,omitempty"`
	Type                        *string        `json:"type,omitempty"`
	Advanced                    bool           `json:"advanced"`
	SelectOptions               []SelectOption `json:"selectOptions,omitempty"`
	SelectOptionsProviderAction *string        `json:"selectOptionsProviderAction,omitempty"`
	Section                     *string        `json:"section,omitempty"`
	Hidden                      *string        `json:"hidden,omitempty"`
	Privacy                     PrivacyLevel   `json:"privacy,omitempty"`
	Placeholder                 *string        `json:"placeholder,omitempty"`
	IsFloat                     bool           `json:"isFloat"`
}

// ProviderMessage is a status message returned by a provider resource.
type ProviderMessage struct {
	Message *string             `json:"message,omitempty"`
	Type    ProviderMessageType `json:"type,omitempty"`
}

// ImportRejectionResource describes why a release was rejected during import.
type ImportRejectionResource struct {
	Reason *string       `json:"reason,omitempty"`
	Type   RejectionType `json:"type,omitempty"`
}

// CustomFormatSpecificationSchema is a specification rule within a custom format.
type CustomFormatSpecificationSchema struct {
	ID                 int                               `json:"id"`
	Name               *string                           `json:"name,omitempty"`
	Implementation     *string                           `json:"implementation,omitempty"`
	ImplementationName *string                           `json:"implementationName,omitempty"`
	InfoLink           *string                           `json:"infoLink,omitempty"`
	Negate             bool                              `json:"negate"`
	Required           bool                              `json:"required"`
	Fields             []Field                           `json:"fields,omitempty"`
	Presets            []CustomFormatSpecificationSchema `json:"presets,omitempty"`
}

// Language represents an audio or subtitle language.
type Language struct {
	ID   int     `json:"id"`
	Name *string `json:"name,omitempty"`
}

// MediaCover is a single artwork image attached to a movie.
type MediaCover struct {
	CoverType MediaCoverType `json:"coverType,omitempty"`
	URL       *string        `json:"url,omitempty"`
	RemoteURL *string        `json:"remoteUrl,omitempty"`
}

// AddMovieOptions controls behaviour when a movie is added to the library.
type AddMovieOptions struct {
	IgnoreEpisodesWithFiles    bool           `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles bool           `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                    MonitorTypes   `json:"monitor,omitempty"`
	SearchForMovie             bool           `json:"searchForMovie,omitempty"`
	AddMethod                  AddMovieMethod `json:"addMethod,omitempty"`
}

// RatingChild holds a single rating value from one source.
type RatingChild struct {
	Votes int        `json:"votes"`
	Value float64    `json:"value"`
	Type  RatingType `json:"type,omitempty"`
}

// Ratings aggregates scores from multiple rating providers.
type Ratings struct {
	IMDB           *RatingChild `json:"imdb,omitempty"`
	TMDB           *RatingChild `json:"tmdb,omitempty"`
	Metacritic     *RatingChild `json:"metacritic,omitempty"`
	RottenTomatoes *RatingChild `json:"rottenTomatoes,omitempty"`
	Trakt          *RatingChild `json:"trakt,omitempty"`
}

// Quality describes the resolution and source of a release.
type Quality struct {
	ID         int           `json:"id"`
	Name       *string       `json:"name,omitempty"`
	Source     QualitySource `json:"source,omitempty"`
	Resolution int           `json:"resolution"`
	Modifier   Modifier      `json:"modifier,omitempty"`
}

// Revision tracks the version and repack status of a release.
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// QualityModel pairs a quality definition with its revision.
type QualityModel struct {
	Quality  Quality  `json:"quality"`
	Revision Revision `json:"revision"`
}

// CustomFormat is a user-defined scoring rule applied to releases.
type CustomFormat struct {
	ID                              int                               `json:"id"`
	Name                            *string                           `json:"name,omitempty"`
	IncludeCustomFormatWhenRenaming *bool                             `json:"includeCustomFormatWhenRenaming,omitempty"`
	Specifications                  []CustomFormatSpecificationSchema `json:"specifications,omitempty"`
}

// MediaInfo holds technical metadata about the video and audio streams.
type MediaInfo struct {
	ID                    int     `json:"id"`
	AudioBitrate          int64   `json:"audioBitrate"`
	AudioChannels         float64 `json:"audioChannels"`
	AudioCodec            *string `json:"audioCodec,omitempty"`
	AudioLanguages        *string `json:"audioLanguages,omitempty"`
	AudioStreamCount      int     `json:"audioStreamCount"`
	VideoBitDepth         int     `json:"videoBitDepth"`
	VideoBitrate          int64   `json:"videoBitrate"`
	VideoCodec            *string `json:"videoCodec,omitempty"`
	VideoFPS              float64 `json:"videoFps"`
	VideoDynamicRange     *string `json:"videoDynamicRange,omitempty"`
	VideoDynamicRangeType *string `json:"videoDynamicRangeType,omitempty"`
	Resolution            *string `json:"resolution,omitempty"`
	RunTime               *string `json:"runTime,omitempty"`
	ScanType              *string `json:"scanType,omitempty"`
	Subtitles             *string `json:"subtitles,omitempty"`
}

// MovieFile represents a physical file on disk associated with a movie.
type MovieFile struct {
	ID                  int            `json:"id"`
	MovieID             int            `json:"movieId"`
	RelativePath        *string        `json:"relativePath,omitempty"`
	Path                *string        `json:"path,omitempty"`
	Size                int64          `json:"size"`
	DateAdded           time.Time      `json:"dateAdded"`
	SceneName           *string        `json:"sceneName,omitempty"`
	ReleaseGroup        *string        `json:"releaseGroup,omitempty"`
	Edition             *string        `json:"edition,omitempty"`
	Languages           []Language     `json:"languages,omitempty"`
	Quality             QualityModel   `json:"quality"`
	CustomFormats       []CustomFormat `json:"customFormats,omitempty"`
	CustomFormatScore   *int           `json:"customFormatScore,omitempty"`
	IndexerFlags        *int           `json:"indexerFlags,omitempty"`
	MediaInfo           *MediaInfo     `json:"mediaInfo,omitempty"`
	OriginalFilePath    *string        `json:"originalFilePath,omitempty"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
}

// MovieCollection groups a movie into a TMDB collection.
type MovieCollection struct {
	Title  *string `json:"title,omitempty"`
	TmdbID int     `json:"tmdbId"`
}

// MovieStatistics holds computed file and size statistics for a movie.
type MovieStatistics struct {
	MovieFileCount int      `json:"movieFileCount"`
	SizeOnDisk     int64    `json:"sizeOnDisk"`
	ReleaseGroups  []string `json:"releaseGroups,omitempty"`
}

// AlternativeTitle is an alternate or translated title for a movie.
type AlternativeTitle struct {
	ID              int        `json:"id"`
	SourceType      SourceType `json:"sourceType,omitempty"`
	MovieMetadataID int        `json:"movieMetadataId"`
	Title           *string    `json:"title,omitempty"`
	CleanTitle      *string    `json:"cleanTitle,omitempty"`
}

// Movie is the central resource; it is referenced by both MovieService and QueueService.
type Movie struct {
	ID                    int                `json:"id"`
	Title                 *string            `json:"title,omitempty"`
	OriginalTitle         *string            `json:"originalTitle,omitempty"`
	OriginalLanguage      Language           `json:"originalLanguage"`
	AlternateTitles       []AlternativeTitle `json:"alternateTitles,omitempty"`
	SecondaryYear         *int               `json:"secondaryYear,omitempty"`
	SecondaryYearSourceID int                `json:"secondaryYearSourceId"`
	SortTitle             *string            `json:"sortTitle,omitempty"`
	SizeOnDisk            *int64             `json:"sizeOnDisk,omitempty"`
	Status                MovieStatusType    `json:"status,omitempty"`
	Overview              *string            `json:"overview,omitempty"`
	InCinemas             *time.Time         `json:"inCinemas,omitempty"`
	PhysicalRelease       *time.Time         `json:"physicalRelease,omitempty"`
	DigitalRelease        *time.Time         `json:"digitalRelease,omitempty"`
	ReleaseDate           *time.Time         `json:"releaseDate,omitempty"`
	PhysicalReleaseNote   *string            `json:"physicalReleaseNote,omitempty"`
	Images                []MediaCover       `json:"images,omitempty"`
	Website               *string            `json:"website,omitempty"`
	RemotePoster          *string            `json:"remotePoster,omitempty"`
	Year                  int                `json:"year"`
	YouTubeTrailerID      *string            `json:"youTubeTrailerId,omitempty"`
	Studio                *string            `json:"studio,omitempty"`
	Path                  *string            `json:"path,omitempty"`
	QualityProfileID      int                `json:"qualityProfileId"`
	HasFile               *bool              `json:"hasFile,omitempty"`
	MovieFileID           int                `json:"movieFileId"`
	Monitored             bool               `json:"monitored"`
	MinimumAvailability   MovieStatusType    `json:"minimumAvailability,omitempty"`
	IsAvailable           bool               `json:"isAvailable"`
	FolderName            *string            `json:"folderName,omitempty"`
	Runtime               int                `json:"runtime"`
	CleanTitle            *string            `json:"cleanTitle,omitempty"`
	ImdbID                *string            `json:"imdbId,omitempty"`
	TmdbID                int                `json:"tmdbId"`
	TitleSlug             *string            `json:"titleSlug,omitempty"`
	RootFolderPath        *string            `json:"rootFolderPath,omitempty"`
	Folder                *string            `json:"folder,omitempty"`
	Certification         *string            `json:"certification,omitempty"`
	Genres                []string           `json:"genres,omitempty"`
	Keywords              []string           `json:"keywords,omitempty"`
	Tags                  []int              `json:"tags,omitempty"`
	Added                 time.Time          `json:"added"`
	AddOptions            *AddMovieOptions   `json:"addOptions,omitempty"`
	Ratings               *Ratings           `json:"ratings,omitempty"`
	MovieFile             *MovieFile         `json:"movieFile,omitempty"`
	Collection            *MovieCollection   `json:"collection,omitempty"`
	Popularity            float64            `json:"popularity"`
	LastSearchTime        *time.Time         `json:"lastSearchTime,omitempty"`
	Statistics            *MovieStatistics   `json:"statistics,omitempty"`
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

// PagedResult is the envelope returned by all paginated Radarr endpoints.
type PagedResult[T any] struct {
	Page         int `json:"page"`
	PageSize     int `json:"pageSize"`
	TotalRecords int `json:"totalRecords"`
	Records      []T `json:"records"`
}

// paginate drives automatic pagination by repeatedly calling getPage with an
// injected page option and yielding each page to the caller via iter.Seq2.
func paginate[T any, O ~func(*resty.Request)](
	ctx context.Context,
	getPage func(context.Context, ...O) (*PagedResult[T], error),
	opts ...O,
) iter.Seq2[PagedResult[T], error] {
	return func(yield func(PagedResult[T], error) bool) {
		page := 1
		for {
			p := page
			pageOpt := O(func(r *resty.Request) {
				r.SetQueryParam("page", strconv.Itoa(p))
			})

			result, err := getPage(ctx, append(slices.Clone(opts), pageOpt)...)
			if err != nil {
				yield(PagedResult[T]{}, err)
				return
			}

			if !yield(*result, nil) {
				return
			}

			if result.PageSize == 0 || page*result.PageSize >= result.TotalRecords {
				return
			}

			page++
		}
	}
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

// derefResult extracts and dereferences a typed value from a resty response.
// resty's SetResult guarantees the type, but we use comma-ok to satisfy the linter.
func derefResult[T any](resp *resty.Response) (T, error) {
	v, ok := resp.Result().(*T)
	if !ok {
		var zero T

		return zero, errors.New("radarr: unexpected response type")
	}

	return *v, nil
}

// ptrResult extracts a typed pointer from a resty response.
func ptrResult[T any](resp *resty.Response) (*T, error) {
	v, ok := resp.Result().(*T)
	if !ok {
		return nil, errors.New("radarr: unexpected response type")
	}

	return v, nil
}
