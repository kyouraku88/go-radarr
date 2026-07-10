package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFilterService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.CustomFilter{
		{ID: 1, Type: new("movieIndex"), Label: new("HD Movies")},
		{ID: 2, Type: new("movieIndex"), Label: new("Wanted")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/customfilter", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFilter.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Label)
	assert.Equal(t, "HD Movies", *got[0].Label)
}

func TestCustomFilterService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.CustomFilter{ID: 3, Type: new("movieIndex"), Label: new("4K")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/customfilter/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFilter.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Label)
	assert.Equal(t, "4K", *got.Label)
}

func TestCustomFilterService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.CustomFilter{Type: new("movieIndex"), Label: new("New Filter")}
	want := radarr.CustomFilter{ID: 4, Type: new("movieIndex"), Label: new("New Filter")}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/customfilter", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFilter.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestCustomFilterService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.CustomFilter{ID: 4, Type: new("movieIndex"), Label: new("Updated Filter")}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/customfilter/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFilter.Update(t.Context(), 4, body)
	require.NoError(t, err)
	require.NotNil(t, got.Label)
	assert.Equal(t, "Updated Filter", *got.Label)
}

func TestCustomFilterService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/customfilter/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.CustomFilter.Delete(t.Context(), 5)
	require.NoError(t, err)
}
