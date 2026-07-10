package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemotePathMappingService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.RemotePathMapping{
		{ID: 1, Host: new("qbittorrent"), RemotePath: new("/downloads"), LocalPath: new("/media/downloads")},
		{ID: 2, Host: new("sabnzbd"), RemotePath: new("/complete"), LocalPath: new("/media/complete")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/remotepathmapping", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RemotePathMapping.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Host)
	assert.Equal(t, "qbittorrent", *got[0].Host)
	require.NotNil(t, got[0].LocalPath)
	assert.Equal(t, "/media/downloads", *got[0].LocalPath)
}

func TestRemotePathMappingService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.RemotePathMapping{ID: 3, Host: new("rtorrent"), RemotePath: new("/torrent"), LocalPath: new("/media/torrent")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/remotepathmapping/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RemotePathMapping.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestRemotePathMappingService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.RemotePathMapping{Host: new("deluge"), RemotePath: new("/done"), LocalPath: new("/media/done")}
	want := radarr.RemotePathMapping{ID: 4, Host: new("deluge"), RemotePath: new("/done"), LocalPath: new("/media/done")}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/remotepathmapping", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RemotePathMapping.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestRemotePathMappingService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.RemotePathMapping{ID: 4, Host: new("deluge"), RemotePath: new("/done"), LocalPath: new("/media/updated")}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/remotepathmapping/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RemotePathMapping.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	require.NotNil(t, got.LocalPath)
	assert.Equal(t, "/media/updated", *got.LocalPath)
}

func TestRemotePathMappingService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/remotepathmapping/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.RemotePathMapping.Delete(t.Context(), 5)
	require.NoError(t, err)
}
