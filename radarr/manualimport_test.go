package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManualImportService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.ManualImportItem{
		{ID: 1, Path: new("/downloads/Movie.2024.mkv"), Name: new("Movie.2024"), Size: 10000000000},
		{ID: 2, Path: new("/downloads/Film.2023.mkv"), Name: new("Film.2023"), Size: 5000000000},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/manualimport", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ManualImport.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Path)
	assert.Equal(t, "/downloads/Movie.2024.mkv", *got[0].Path)
}

func TestManualImportService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/manualimport", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "/downloads", q.Get("folder"))
		assert.Equal(t, "abc123", q.Get("downloadId"))
		assert.Equal(t, "10", q.Get("movieId"))
		assert.Equal(t, "true", q.Get("filterExistingFiles"))
		writeJSON(w, http.StatusOK, []radarr.ManualImportItem{})
	})

	c := newTestClient(t, mux)
	_, err := c.ManualImport.List(t.Context(),
		radarr.WithManualImportFolder("/downloads"),
		radarr.WithManualImportDownloadID("abc123"),
		radarr.WithManualImportMovieID(10),
		radarr.WithManualImportFilterExistingFiles(true),
	)
	require.NoError(t, err)
}

func TestManualImportService_Reprocess(t *testing.T) {
	t.Parallel()

	body := []radarr.ManualImportReprocessItem{
		{ID: 1, MovieID: 42, Path: new("/downloads/Movie.2024.mkv")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/manualimport", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ManualImport.Reprocess(t.Context(), body)
	require.NoError(t, err)
}
