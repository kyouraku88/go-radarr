package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.RenameMovie{
		{ID: 1, MovieID: 10, MovieFileID: 100, ExistingPath: new("/movies/Fight Club (1999)/Fight Club (1999).mkv"), NewPath: new("/movies/Fight Club (1999)/Fight.Club.1999.1080p.mkv")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/rename", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Rename.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 10, got[0].MovieID)
	require.NotNil(t, got[0].ExistingPath)
	assert.Equal(t, "/movies/Fight Club (1999)/Fight Club (1999).mkv", *got[0].ExistingPath)
}

func TestRenameService_List_WithMovieID(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/rename", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("movieId"))
		writeJSON(w, http.StatusOK, []radarr.RenameMovie{})
	})

	c := newTestClient(t, mux)
	_, err := c.Rename.List(t.Context(), radarr.WithRenameMovieID(10))
	require.NoError(t, err)
}
