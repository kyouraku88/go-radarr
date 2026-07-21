package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadClientService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.DownloadClient{
		{ID: 1, Name: new("qBittorrent"), Enable: true, Priority: 1},
		{ID: 2, Name: new("SABnzbd"), Enable: false, Priority: 2},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/downloadclient", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "qBittorrent", *got[0].Name)
	assert.True(t, got[0].Enable)
}

func TestDownloadClientService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.DownloadClient{ID: 3, Name: new("Transmission"), Enable: true, Protocol: radarr.DownloadProtocolTorrent}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/downloadclient/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.Equal(t, radarr.DownloadProtocolTorrent, got.Protocol)
}

func TestDownloadClientService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.DownloadClient{{ID: 0, Name: new("qBittorrent")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/downloadclient/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestDownloadClientService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.DownloadClient{Name: new("NZBGet"), Enable: true}
	want := radarr.DownloadClient{ID: 4, Name: new("NZBGet"), Enable: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/downloadclient", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestDownloadClientService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.DownloadClient{ID: 4, Name: new("NZBGet"), Enable: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/downloadclient/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.Enable)
}

func TestDownloadClientService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := radarr.DownloadClientBulkRequest{IDs: []int{1, 2}, Enable: new(true)}
	want := []radarr.DownloadClient{
		{ID: 1, Name: new("qBittorrent"), Enable: true},
		{ID: 2, Name: new("SABnzbd"), Enable: true},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/downloadclient/bulk", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.DownloadClient.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.True(t, got[0].Enable)
}

func TestDownloadClientService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/downloadclient/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DownloadClient.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestDownloadClientService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/downloadclient/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DownloadClient.DeleteBulk(t.Context(), radarr.DownloadClientBulkRequest{IDs: []int{1, 2}})
	require.NoError(t, err)
}

func TestDownloadClientService_Test(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/downloadclient/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DownloadClient.Test(t.Context(), radarr.DownloadClient{Name: new("qBittorrent")})
	require.NoError(t, err)
}

func TestDownloadClientService_TestAll(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/downloadclient/testall", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DownloadClient.TestAll(t.Context())
	require.NoError(t, err)
}

func TestDownloadClientService_Action(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/downloadclient/action/checkHealth", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.DownloadClient.Action(t.Context(), "checkHealth", radarr.DownloadClient{Name: new("qBittorrent")})
	require.NoError(t, err)
}
