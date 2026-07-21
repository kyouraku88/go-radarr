package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectionService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Collection{
		{ID: 1, TmdbID: 10, Title: new("The Dark Knight Collection"), MissingMovies: 1},
		{ID: 2, TmdbID: 20, Title: new("Marvel Cinematic Universe"), MissingMovies: 0},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/collection", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Collection.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "The Dark Knight Collection", *got[0].Title)
	assert.Equal(t, 1, got[0].MissingMovies)
}

func TestCollectionService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Collection{ID: 3, TmdbID: 30, Title: new("Alien Collection"), Monitored: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/collection/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Collection.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Alien Collection", *got.Title)
	assert.True(t, got.Monitored)
}

func TestCollectionService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.Collection{ID: 3, TmdbID: 30, Title: new("Alien Collection"), Monitored: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/collection/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Collection.Update(t.Context(), 3, body)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.False(t, got.Monitored)
}

func TestCollectionService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := radarr.CollectionUpdateRequest{
		CollectionIDs: []int{1, 2},
		Monitored:     new(true),
	}
	want := []radarr.Collection{
		{ID: 1, TmdbID: 10, Monitored: true},
		{ID: 2, TmdbID: 20, Monitored: true},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/collection", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Collection.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.True(t, got[0].Monitored)
}
