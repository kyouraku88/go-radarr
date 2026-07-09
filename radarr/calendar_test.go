package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 1, TmdbID: 100}, {ID: 2, TmdbID: 200}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/calendar", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Calendar.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, 1, got[0].ID)
	assert.Equal(t, 200, got[1].TmdbID)
}

func TestCalendarService_List_Options(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/calendar", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "2024-01-01T00:00:00Z", q.Get("start"))
		assert.Equal(t, "2024-01-31T00:00:00Z", q.Get("end"))
		assert.Equal(t, "true", q.Get("unmonitored"))
		assert.Equal(t, "action,drama", q.Get("tags"))
		writeJSON(w, http.StatusOK, []radarr.Movie{})
	})

	c := newTestClient(t, mux)
	_, err := c.Calendar.List(t.Context(),
		radarr.WithCalendarStart(start),
		radarr.WithCalendarEnd(end),
		radarr.WithCalendarUnmonitored(true),
		radarr.WithCalendarTags("action,drama"),
	)
	require.NoError(t, err)
}
