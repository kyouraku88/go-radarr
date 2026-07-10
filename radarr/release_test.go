package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReleaseService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Release{
		{ID: 1, Title: new("Movie.2024.1080p.BluRay"), Size: 10000000000, Approved: true, PublishDate: time.Now()},
		{ID: 2, Title: new("Movie.2024.720p.WEB-DL"), Size: 5000000000, Approved: false, PublishDate: time.Now()},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/release", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Release.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "Movie.2024.1080p.BluRay", *got[0].Title)
	assert.True(t, got[0].Approved)
}

func TestReleaseService_List_WithMovieID(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/release", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "99", r.URL.Query().Get("movieId"))
		writeJSON(w, http.StatusOK, []radarr.Release{})
	})

	c := newTestClient(t, mux)
	_, err := c.Release.List(t.Context(), radarr.WithReleaseMovieID(99))
	require.NoError(t, err)
}

func TestReleaseService_Download(t *testing.T) {
	t.Parallel()

	body := radarr.Release{GUID: new("nzb-guid-123"), Title: new("Movie.2024.1080p")}
	want := radarr.Release{ID: 1, GUID: new("nzb-guid-123"), Title: new("Movie.2024.1080p")}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/release", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Release.Download(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
	require.NotNil(t, got.GUID)
	assert.Equal(t, "nzb-guid-123", *got.GUID)
}

func TestReleaseService_Push(t *testing.T) {
	t.Parallel()

	body := radarr.Release{Title: new("Movie.2024.1080p"), DownloadURL: new("http://example.com/release.nzb")}
	want := []radarr.Release{
		{ID: 1, Title: new("Movie.2024.1080p"), Approved: true},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/release/push", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Release.Push(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].Approved)
}
