package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootFolderService_List(t *testing.T) {
	t.Parallel()

	freeSpace := int64(500_000_000_000)
	want := []radarr.RootFolder{
		{ID: 1, Path: new("/mnt/media"), Accessible: true, FreeSpace: &freeSpace},
		{ID: 2, Path: new("/mnt/backup"), Accessible: false},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/rootfolder", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RootFolder.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Path)
	assert.Equal(t, "/mnt/media", *got[0].Path)
	assert.True(t, got[0].Accessible)
	require.NotNil(t, got[0].FreeSpace)
	assert.Equal(t, int64(500_000_000_000), *got[0].FreeSpace)
	assert.False(t, got[1].Accessible)
}

func TestRootFolderService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.RootFolder{
		ID:         3,
		Path:       new("/data"),
		Accessible: true,
		UnmappedFolders: []radarr.UnmappedFolder{
			{Name: new("orphan"), Path: new("/data/orphan")},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/rootfolder/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.RootFolder.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.Len(t, got.UnmappedFolders, 1)
	require.NotNil(t, got.UnmappedFolders[0].Name)
	assert.Equal(t, "orphan", *got.UnmappedFolders[0].Name)
}

func TestRootFolderService_Create(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/rootfolder", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.RootFolder
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotNil(t, body.Path)
		assert.Equal(t, "/new/path", *body.Path)
		body.ID = 5
		writeJSON(w, http.StatusCreated, body)
	})

	c := newTestClient(t, mux)
	got, err := c.RootFolder.Create(t.Context(), radarr.RootFolder{Path: new("/new/path")})
	require.NoError(t, err)
	assert.Equal(t, 5, got.ID)
}

func TestRootFolderService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/rootfolder/4", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.RootFolder.Delete(t.Context(), 4)
	require.NoError(t, err)
}
