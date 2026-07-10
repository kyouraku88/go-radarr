package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportListService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.ImportList{
		{ID: 1, Name: new("TMDb Popular"), Enabled: true, EnableAuto: true, ListOrder: 1},
		{ID: 2, Name: new("Trakt Trending"), Enabled: false, EnableAuto: false, ListOrder: 2},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/importlist", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "TMDb Popular", *got[0].Name)
	assert.True(t, got[0].Enabled)
}

func TestImportListService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.ImportList{ID: 3, Name: new("IMDb List"), Enabled: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/importlist/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestImportListService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.ImportList{{ID: 0, Name: new("TMDbPopularImport")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/importlist/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestImportListService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.ImportList{Name: new("New List"), Enabled: true}
	want := radarr.ImportList{ID: 4, Name: new("New List"), Enabled: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/importlist", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestImportListService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.ImportList{ID: 4, Name: new("Updated List"), Enabled: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/importlist/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.Enabled)
}

func TestImportListService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := radarr.ImportListBulkRequest{IDs: []int{1, 2}, Enabled: new(true)}
	want := []radarr.ImportList{{ID: 1, Enabled: true}, {ID: 2, Enabled: true}}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/importlist/bulk", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestImportListService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/importlist/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ImportList.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestImportListService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/importlist/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ImportList.DeleteBulk(t.Context(), radarr.ImportListBulkRequest{IDs: []int{1, 2}})
	require.NoError(t, err)
}

func TestImportListService_Test(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/importlist/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ImportList.Test(t.Context(), radarr.ImportList{Name: new("Test List")})
	require.NoError(t, err)
}

func TestImportListService_TestAll(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/importlist/testall", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ImportList.TestAll(t.Context())
	require.NoError(t, err)
}

func TestImportListService_Action(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/importlist/action/getMovies", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.ImportList.Action(t.Context(), "getMovies", radarr.ImportList{Name: new("TMDb Popular")})
	require.NoError(t, err)
}

func TestImportListService_Movies(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{
		{ID: 1, Title: new("Interstellar"), TmdbID: 157336},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/importlist/movie", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "true", r.URL.Query().Get("includeRecommendations"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.ImportList.Movies(t.Context(), radarr.WithImportListMoviesIncludeRecommendations(true))
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "Interstellar", *got[0].Title)
}
