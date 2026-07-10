package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexerService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Indexer{
		{ID: 1, Name: new("Prowlarr"), EnableRss: true, EnableAutomaticSearch: true, Priority: 25},
		{ID: 2, Name: new("Jackett"), EnableRss: false, EnableAutomaticSearch: true, Priority: 50},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/indexer", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Prowlarr", *got[0].Name)
	assert.True(t, got[0].EnableRss)
}

func TestIndexerService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Indexer{ID: 3, Name: new("NZBGeek"), EnableRss: true, Priority: 25}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/indexer/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestIndexerService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.Indexer{{ID: 0, Name: new("Newznab")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/indexer/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestIndexerService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.Indexer{Name: new("New Indexer"), EnableRss: true}
	want := radarr.Indexer{ID: 4, Name: new("New Indexer"), EnableRss: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/indexer", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestIndexerService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.Indexer{ID: 4, Name: new("New Indexer"), EnableRss: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/indexer/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.EnableRss)
}

func TestIndexerService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := radarr.IndexerBulkRequest{IDs: []int{1, 2}, EnableRss: new(true)}
	want := []radarr.Indexer{{ID: 1, EnableRss: true}, {ID: 2, EnableRss: true}}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/indexer/bulk", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Indexer.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.True(t, got[0].EnableRss)
}

func TestIndexerService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/indexer/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Indexer.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestIndexerService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/indexer/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Indexer.DeleteBulk(t.Context(), radarr.IndexerBulkRequest{IDs: []int{1, 2}})
	require.NoError(t, err)
}

func TestIndexerService_Test(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/indexer/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Indexer.Test(t.Context(), radarr.Indexer{Name: new("Prowlarr")})
	require.NoError(t, err)
}

func TestIndexerService_TestAll(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/indexer/testall", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Indexer.TestAll(t.Context())
	require.NoError(t, err)
}

func TestIndexerService_Action(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/indexer/action/checkCaps", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Indexer.Action(t.Context(), "checkCaps", radarr.Indexer{Name: new("Prowlarr")})
	require.NoError(t, err)
}
