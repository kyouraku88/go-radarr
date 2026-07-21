package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelayProfileService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.DelayProfile{
		{ID: 1, EnableUsenet: true, EnableTorrent: true, UsenetDelay: 0, TorrentDelay: 60},
		{ID: 2, EnableUsenet: false, EnableTorrent: true, TorrentDelay: 120},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/delayprofile", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DelayProfile.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.True(t, got[0].EnableUsenet)
	assert.Equal(t, 60, got[0].TorrentDelay)
}

func TestDelayProfileService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.DelayProfile{ID: 3, EnableUsenet: true, EnableTorrent: false, UsenetDelay: 30}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/delayprofile/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DelayProfile.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.Equal(t, 30, got.UsenetDelay)
}

func TestDelayProfileService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.DelayProfile{EnableUsenet: true, EnableTorrent: true}
	want := radarr.DelayProfile{ID: 4, EnableUsenet: true, EnableTorrent: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/delayprofile", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DelayProfile.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestDelayProfileService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.DelayProfile{ID: 4, EnableUsenet: true, EnableTorrent: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/delayprofile/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DelayProfile.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.EnableTorrent)
}

func TestDelayProfileService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/delayprofile/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DelayProfile.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestDelayProfileService_Reorder(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/delayprofile/reorder/3", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "1", r.URL.Query().Get("after"))
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DelayProfile.Reorder(t.Context(), 3, radarr.WithReorderAfter(1))
	require.NoError(t, err)
}
