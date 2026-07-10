package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReleaseProfileService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.ReleaseProfile{
		{ID: 1, Name: new("Preferred Terms"), Enabled: true, IndexerID: 0},
		{ID: 2, Name: new("Ignored Terms"), Enabled: false, IndexerID: 1},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/releaseprofile", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ReleaseProfile.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Preferred Terms", *got[0].Name)
	assert.True(t, got[0].Enabled)
}

func TestReleaseProfileService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.ReleaseProfile{ID: 3, Name: new("Custom Profile"), Enabled: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/releaseprofile/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ReleaseProfile.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestReleaseProfileService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.ReleaseProfile{Name: new("New Profile"), Enabled: true}
	want := radarr.ReleaseProfile{ID: 4, Name: new("New Profile"), Enabled: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/releaseprofile", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ReleaseProfile.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestReleaseProfileService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.ReleaseProfile{ID: 4, Name: new("Updated Profile"), Enabled: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/releaseprofile/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ReleaseProfile.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.Enabled)
}

func TestReleaseProfileService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/releaseprofile/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ReleaseProfile.Delete(t.Context(), 5)
	require.NoError(t, err)
}
