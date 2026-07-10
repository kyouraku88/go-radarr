package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.MetadataProvider{
		{ID: 1, Name: new("Kodi (XBMC) / Emby"), Enable: true},
		{ID: 2, Name: new("Plex"), Enable: false},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/metadata", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Metadata.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Kodi (XBMC) / Emby", *got[0].Name)
	assert.True(t, got[0].Enable)
}

func TestMetadataService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.MetadataProvider{ID: 3, Name: new("Roksbox"), Enable: false}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/metadata/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Metadata.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestMetadataService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.MetadataProvider{{ID: 0, Name: new("KodiMetadata")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/metadata/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Metadata.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestMetadataService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.MetadataProvider{Name: new("New Agent"), Enable: true}
	want := radarr.MetadataProvider{ID: 4, Name: new("New Agent"), Enable: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/metadata", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Metadata.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestMetadataService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.MetadataProvider{ID: 4, Name: new("Updated Agent"), Enable: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/metadata/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Metadata.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.Enable)
}

func TestMetadataService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/metadata/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Metadata.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestMetadataService_Test(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/metadata/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Metadata.Test(t.Context(), radarr.MetadataProvider{Name: new("Kodi")})
	require.NoError(t, err)
}

func TestMetadataService_TestAll(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/metadata/testall", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Metadata.TestAll(t.Context())
	require.NoError(t, err)
}

func TestMetadataService_Action(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/metadata/action/exportNfo", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Metadata.Action(t.Context(), "exportNfo", radarr.MetadataProvider{Name: new("Kodi")})
	require.NoError(t, err)
}
