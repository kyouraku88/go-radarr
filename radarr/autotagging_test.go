package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutoTaggingService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.AutoTagging{
		{ID: 1, Name: new("Tag Action"), RemoveTagsAutomatically: true},
		{ID: 2, Name: new("Tag HD"), RemoveTagsAutomatically: false},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/autotagging", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AutoTagging.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Tag Action", *got[0].Name)
	assert.True(t, got[0].RemoveTagsAutomatically)
}

func TestAutoTaggingService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.AutoTagging{ID: 3, Name: new("Tag Thriller")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/autotagging/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AutoTagging.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Name)
	assert.Equal(t, "Tag Thriller", *got.Name)
}

func TestAutoTaggingService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.AutoTagging{{ID: 0, Name: new("GenreSpecification")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/autotagging/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AutoTagging.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestAutoTaggingService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.AutoTagging{Name: new("Tag Comedy"), RemoveTagsAutomatically: false}
	want := radarr.AutoTagging{ID: 4, Name: new("Tag Comedy"), RemoveTagsAutomatically: false}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/autotagging", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AutoTagging.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestAutoTaggingService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.AutoTagging{ID: 4, Name: new("Tag Comedy Updated")}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/autotagging/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.AutoTagging.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestAutoTaggingService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/autotagging/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.AutoTagging.Delete(t.Context(), 5)
	require.NoError(t, err)
}
