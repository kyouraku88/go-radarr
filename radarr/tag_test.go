package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Tag{{ID: 1, Label: new("4k")}, {ID: 2, Label: new("hdr")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/tag", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Label)
	assert.Equal(t, "4k", *got[0].Label)
}

func TestTagService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Tag{ID: 3, Label: new("remux")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/tag/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Label)
	assert.Equal(t, "remux", *got.Label)
}

func TestTagService_Create(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/tag", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.Tag
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotNil(t, body.Label)
		assert.Equal(t, "new-tag", *body.Label)
		body.ID = 10
		writeJSON(w, http.StatusCreated, body)
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.Create(t.Context(), radarr.Tag{Label: new("new-tag")})
	require.NoError(t, err)
	assert.Equal(t, 10, got.ID)
}

func TestTagService_Update(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/tag/5", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.Tag
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotNil(t, body.Label)
		assert.Equal(t, "updated", *body.Label)
		writeJSON(w, http.StatusOK, radarr.Tag{ID: 5, Label: body.Label})
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.Update(t.Context(), 5, radarr.Tag{ID: 5, Label: new("updated")})
	require.NoError(t, err)
	assert.Equal(t, 5, got.ID)
}

func TestTagService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/tag/7", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Tag.Delete(t.Context(), 7)
	require.NoError(t, err)
}

func TestTagService_ListDetails(t *testing.T) {
	t.Parallel()

	want := []radarr.TagDetails{
		{ID: 1, Label: new("4k"), MovieIDs: []int{10, 20, 30}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/tag/detail", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.ListDetails(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, []int{10, 20, 30}, got[0].MovieIDs)
}

func TestTagService_GetDetail(t *testing.T) {
	t.Parallel()

	want := radarr.TagDetails{
		ID:                2,
		Label:             new("hdr"),
		MovieIDs:          []int{5},
		DownloadClientIDs: []int{1},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/tag/detail/2", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Tag.GetDetail(t.Context(), 2)
	require.NoError(t, err)
	assert.Equal(t, 2, got.ID)
	assert.Equal(t, []int{5}, got.MovieIDs)
	assert.Equal(t, []int{1}, got.DownloadClientIDs)
}
