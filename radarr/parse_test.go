package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseService_Parse(t *testing.T) {
	t.Parallel()

	want := radarr.ParseResult{
		ID:    1,
		Title: new("Fight.Club.1999.1080p.BluRay"),
		ParsedMovieInfo: &radarr.ParsedMovieInfo{
			MovieTitle: new("Fight Club"),
			Year:       1999,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/parse", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "Fight.Club.1999.1080p.BluRay", r.URL.Query().Get("title"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Parse.Parse(t.Context(), radarr.WithParseTitle("Fight.Club.1999.1080p.BluRay"))
	require.NoError(t, err)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Fight.Club.1999.1080p.BluRay", *got.Title)
	require.NotNil(t, got.ParsedMovieInfo)
	require.NotNil(t, got.ParsedMovieInfo.MovieTitle)
	assert.Equal(t, "Fight Club", *got.ParsedMovieInfo.MovieTitle)
	assert.Equal(t, 1999, got.ParsedMovieInfo.Year)
}

func TestParseService_Parse_NoOptions(t *testing.T) {
	t.Parallel()

	want := radarr.ParseResult{ID: 2}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/parse", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Parse.Parse(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, got.ID)
}
