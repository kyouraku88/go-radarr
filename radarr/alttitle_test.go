package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAltTitleService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.AlternativeTitle{
		{ID: 1, MovieMetadataID: 10, Title: new("Fight Club")},
		{ID: 2, MovieMetadataID: 10, Title: new("Club de la Lutte")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/alttitle", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AltTitle.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "Fight Club", *got[0].Title)
}

func TestAltTitleService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/alttitle", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "42", q.Get("movieId"))
		assert.Equal(t, "7", q.Get("movieMetadataId"))
		writeJSON(w, http.StatusOK, []radarr.AlternativeTitle{})
	})

	c := newTestClient(t, mux)
	_, err := c.AltTitle.List(t.Context(),
		radarr.WithAltTitleMovieID(42),
		radarr.WithAltTitleMovieMetadataID(7),
	)
	require.NoError(t, err)
}

func TestAltTitleService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.AlternativeTitle{ID: 5, MovieMetadataID: 20, Title: new("Pulp Fiction")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/alttitle/5", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AltTitle.Get(t.Context(), 5)
	require.NoError(t, err)
	assert.Equal(t, 5, got.ID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Pulp Fiction", *got.Title)
}
