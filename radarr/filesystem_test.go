package radarr_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilesystemService_Browse(t *testing.T) {
	t.Parallel()

	want := radarr.FilesystemResult{
		Parent: new("/media"),
		Directories: []radarr.FilesystemItem{
			{Name: new("movies"), Path: new("/media/movies"), Type: new("folder")},
		},
		Files: []radarr.FilesystemItem{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/filesystem", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "/media", r.URL.Query().Get("path"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Filesystem.Browse(t.Context(), radarr.WithFilesystemPath("/media"))
	require.NoError(t, err)
	require.NotNil(t, got.Parent)
	assert.Equal(t, "/media", *got.Parent)
	require.Len(t, got.Directories, 1)
	require.NotNil(t, got.Directories[0].Name)
	assert.Equal(t, "movies", *got.Directories[0].Name)
}

func TestFilesystemService_Browse_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/filesystem", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "true", q.Get("includeFiles"))
		assert.Equal(t, "true", q.Get("allowFoldersWithoutTrailingSlashes"))
		writeJSON(w, http.StatusOK, radarr.FilesystemResult{})
	})

	c := newTestClient(t, mux)
	_, err := c.Filesystem.Browse(t.Context(),
		radarr.WithFilesystemIncludeFiles(true),
		radarr.WithFilesystemAllowFoldersWithoutTrailingSlashes(true),
	)
	require.NoError(t, err)
}

func TestFilesystemService_MediaFiles(t *testing.T) {
	t.Parallel()

	want := []radarr.FilesystemItem{
		{Name: new("Movie.2024.mkv"), Path: new("/media/movies/Movie.2024.mkv")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/filesystem/mediafiles", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "/media/movies", r.URL.Query().Get("path"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Filesystem.MediaFiles(t.Context(), radarr.WithMediaFilesPath("/media/movies"))
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Movie.2024.mkv", *got[0].Name)
}

func TestFilesystemService_FileType(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/filesystem/type", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "/media/movies", r.URL.Query().Get("path"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `"folder"`)
	})

	c := newTestClient(t, mux)
	got, err := c.Filesystem.FileType(t.Context(), radarr.WithFileTypePath("/media/movies"))
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "folder", *got)
}
