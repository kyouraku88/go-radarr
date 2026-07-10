package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtraFileService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.ExtraFile{
		{ID: 1, MovieID: 10, RelativePath: new("Movie.2024.srt"), Extension: new(".srt"), Type: radarr.ExtraFileTypeSubtitle},
		{ID: 2, MovieID: 10, RelativePath: new("Movie.2024.nfo"), Extension: new(".nfo"), Type: radarr.ExtraFileTypeOther},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/extrafile", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ExtraFile.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].RelativePath)
	assert.Equal(t, "Movie.2024.srt", *got[0].RelativePath)
	assert.Equal(t, radarr.ExtraFileTypeSubtitle, got[0].Type)
}

func TestExtraFileService_List_WithMovieID(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/extrafile", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "42", r.URL.Query().Get("movieId"))
		writeJSON(w, http.StatusOK, []radarr.ExtraFile{})
	})

	c := newTestClient(t, mux)
	_, err := c.ExtraFile.List(t.Context(), radarr.WithExtraFileMovieID(42))
	require.NoError(t, err)
}
