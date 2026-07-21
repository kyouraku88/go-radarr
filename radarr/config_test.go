package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigService_GetHostConfig(t *testing.T) {
	t.Parallel()

	want := radarr.HostConfig{ID: 1, Port: 7878, EnableSsl: false, AnalyticsEnabled: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/host", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetHostConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.Equal(t, 7878, got.Port)
	assert.False(t, got.EnableSsl)
}

func TestConfigService_UpdateHostConfig(t *testing.T) {
	t.Parallel()

	body := radarr.HostConfig{ID: 1, Port: 7878, EnableSsl: true}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/config/host/1", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.UpdateHostConfig(t.Context(), 1, body)
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.True(t, got.EnableSsl)
}

func TestConfigService_GetUIConfig(t *testing.T) {
	t.Parallel()

	want := radarr.UIConfig{ID: 1, CalendarWeekColumnHeader: new("ddd M/D")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/ui", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetUIConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	require.NotNil(t, got.CalendarWeekColumnHeader)
	assert.Equal(t, "ddd M/D", *got.CalendarWeekColumnHeader)
}

func TestConfigService_UpdateUIConfig(t *testing.T) {
	t.Parallel()

	body := radarr.UIConfig{ID: 1, CalendarWeekColumnHeader: new("ddd D/M")}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/config/ui/1", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.UpdateUIConfig(t.Context(), 1, body)
	require.NoError(t, err)
	require.NotNil(t, got.CalendarWeekColumnHeader)
	assert.Equal(t, "ddd D/M", *got.CalendarWeekColumnHeader)
}

func TestConfigService_GetNamingConfig(t *testing.T) {
	t.Parallel()

	want := radarr.NamingConfig{ID: 1, RenameMovies: true, ReplaceIllegalCharacters: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/naming", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetNamingConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.True(t, got.RenameMovies)
}

func TestConfigService_UpdateNamingConfig(t *testing.T) {
	t.Parallel()

	body := radarr.NamingConfig{ID: 1, RenameMovies: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/config/naming/1", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.UpdateNamingConfig(t.Context(), 1, body)
	require.NoError(t, err)
	assert.False(t, got.RenameMovies)
}

func TestConfigService_GetNamingExamples(t *testing.T) {
	t.Parallel()

	want := map[string]string{
		"singleEpisodeExample": "Fight Club (1999).mkv",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/naming/examples", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "true", r.URL.Query().Get("renameMovies"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetNamingExamples(t.Context(), radarr.WithNamingExamplesRenameMovies(true))
	require.NoError(t, err)
	assert.Equal(t, "Fight Club (1999).mkv", got["singleEpisodeExample"])
}

func TestConfigService_GetMediaManagementConfig(t *testing.T) {
	t.Parallel()

	want := radarr.MediaManagementConfig{ID: 1, CreateEmptyMovieFolders: true, DeleteEmptyFolders: false}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/mediamanagement", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetMediaManagementConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.True(t, got.CreateEmptyMovieFolders)
}

func TestConfigService_GetDownloadClientConfig(t *testing.T) {
	t.Parallel()

	want := radarr.DownloadClientConfig{ID: 1, DownloadClientWorkingFolders: new("_UNPACK_,_FAILED_")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/downloadclient", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetDownloadClientConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
}

func TestConfigService_GetIndexerConfig(t *testing.T) {
	t.Parallel()

	want := radarr.IndexerConfig{ID: 1, MinimumAge: 0, Retention: 0, MaximumSize: 0, RssSyncInterval: 60}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/indexer", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetIndexerConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.Equal(t, 60, got.RssSyncInterval)
}

func TestConfigService_GetMetadataConfig(t *testing.T) {
	t.Parallel()

	want := radarr.MetadataConfig{ID: 1, CertificationCountry: radarr.TMDbCountryCode("us")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/config/metadata", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Config.GetMetadataConfig(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	assert.Equal(t, radarr.TMDbCountryCode("us"), got.CertificationCountry)
}
