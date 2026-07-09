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

	Movie          *MovieService
	MovieEditor    *MovieEditorService
	MovieFile      *MovieFileService
	Queue          *QueueService
	History        *HistoryService
	Calendar       *CalendarService
	System         *SystemService
	Tag            *TagService
	QualityProfile *QualityProfileService
	RootFolder     *RootFolderService
	Command        *CommandService
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
	r.Movie = (*MovieService)(&r.common)
	r.MovieEditor = (*MovieEditorService)(&r.common)
	r.MovieFile = (*MovieFileService)(&r.common)
	r.Queue = (*QueueService)(&r.common)
	r.History = (*HistoryService)(&r.common)
	r.Calendar = (*CalendarService)(&r.common)
	r.System = (*SystemService)(&r.common)
	r.Tag = (*TagService)(&r.common)
	r.QualityProfile = (*QualityProfileService)(&r.common)
	r.RootFolder = (*RootFolderService)(&r.common)
	r.Command = (*CommandService)(&r.common)

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

// ---------------------------------------------------------------------------
// Shared value types
// ---------------------------------------------------------------------------

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
	ID                              int     `json:"id"`
	Name                            *string `json:"name,omitempty"`
	IncludeCustomFormatWhenRenaming *bool   `json:"includeCustomFormatWhenRenaming,omitempty"`
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
