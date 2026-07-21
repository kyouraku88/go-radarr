package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieFileService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.MovieFile{
		{ID: 1, MovieID: 10, Size: 1024, DateAdded: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, MovieID: 20, Size: 2048, DateAdded: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/moviefile", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.MovieFile.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, 1, got[0].ID)
	assert.Equal(t, 2048, int(got[1].Size))
}

func TestMovieFileService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/moviefile", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, []string{"10", "11"}, r.URL.Query()["movieId"])
		assert.Equal(t, []string{"1", "2", "3"}, r.URL.Query()["movieFileIds"])
		writeJSON(w, http.StatusOK, []radarr.MovieFile{})
	})

	c := newTestClient(t, mux)
	_, err := c.MovieFile.List(t.Context(),
		radarr.WithMovieFileMovieIDs(10, 11),
		radarr.WithMovieFileIDs(1, 2, 3),
	)
	require.NoError(t, err)
}

func TestMovieFileService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.MovieFile{ID: 7, MovieID: 42, Size: 4096, DateAdded: time.Now().UTC().Truncate(time.Second)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/moviefile/7", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.MovieFile.Get(t.Context(), 7)
	require.NoError(t, err)
	assert.Equal(t, 7, got.ID)
	assert.Equal(t, 42, got.MovieID)
}

func TestMovieFileService_Update(t *testing.T) {
	t.Parallel()

	input := radarr.MovieFile{ID: 3, MovieID: 10, Size: 512, DateAdded: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/moviefile/3", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.MovieFile
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, 3, body.ID)
		writeJSON(w, http.StatusOK, body)
	})

	c := newTestClient(t, mux)
	got, err := c.MovieFile.Update(t.Context(), 3, input)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestMovieFileService_UpdateBulk(t *testing.T) {
	t.Parallel()

	files := []radarr.MovieFile{
		{ID: 1, MovieID: 10, DateAdded: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, MovieID: 10, DateAdded: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/moviefile/bulk", func(w http.ResponseWriter, r *http.Request) {
		var body []radarr.MovieFile
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Len(t, body, 2)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.MovieFile.UpdateBulk(t.Context(), files)
	require.NoError(t, err)
}

func TestMovieFileService_UpdateEditor(t *testing.T) {
	t.Parallel()

	body := radarr.MovieFileBulkRequest{
		MovieFileIDs: []int{1, 2, 3},
		Edition:      new("Director's Cut"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/moviefile/editor", func(w http.ResponseWriter, r *http.Request) {
		var got radarr.MovieFileBulkRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, []int{1, 2, 3}, got.MovieFileIDs)
		assert.NotNil(t, got.Edition)
		assert.Equal(t, "Director's Cut", *got.Edition)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.MovieFile.UpdateEditor(t.Context(), body)
	require.NoError(t, err)
}

func TestMovieFileService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/moviefile/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.MovieFile.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestMovieFileService_DeleteBulk(t *testing.T) {
	t.Parallel()

	req := radarr.MovieFileBulkRequest{MovieFileIDs: []int{10, 11, 12}}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/moviefile/bulk", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.MovieFileBulkRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{10, 11, 12}, body.MovieFileIDs)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.MovieFile.DeleteBulk(t.Context(), req)
	require.NoError(t, err)
}
