package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocklistService_List(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.BlocklistRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 2,
		Records: []radarr.BlocklistRecord{
			{ID: 1, MovieID: 100, SourceTitle: new("Movie.2024.1080p"), Date: time.Now()},
			{ID: 2, MovieID: 200, SourceTitle: new("Movie.2023.720p"), Date: time.Now()},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/blocklist", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Blocklist.List(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, got.TotalRecords)
	require.Len(t, got.Records, 2)
	require.NotNil(t, got.Records[0].SourceTitle)
	assert.Equal(t, "Movie.2024.1080p", *got.Records[0].SourceTitle)
}

func TestBlocklistService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/blocklist", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "25", q.Get("pageSize"))
		assert.Equal(t, "date", q.Get("sortKey"))
		assert.Equal(t, "descending", q.Get("sortDirection"))
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.BlocklistRecord]{})
	})

	c := newTestClient(t, mux)
	_, err := c.Blocklist.List(t.Context(),
		radarr.WithBlocklistPage(2),
		radarr.WithBlocklistPageSize(25),
		radarr.WithBlocklistSortKey("date"),
		radarr.WithBlocklistSortDirection(radarr.SortDirectionDescending),
	)
	require.NoError(t, err)
}

func TestBlocklistService_ListWithPagination(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.BlocklistRecord]{
		{Page: 1, PageSize: 2, TotalRecords: 5, Records: []radarr.BlocklistRecord{{ID: 1, MovieID: 10}, {ID: 2, MovieID: 20}}},
		{Page: 2, PageSize: 2, TotalRecords: 5, Records: []radarr.BlocklistRecord{{ID: 3, MovieID: 30}, {ID: 4, MovieID: 40}}},
		{Page: 3, PageSize: 2, TotalRecords: 5, Records: []radarr.BlocklistRecord{{ID: 5, MovieID: 50}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/blocklist", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, pages[0])
		case "2":
			writeJSON(w, http.StatusOK, pages[1])
		case "3":
			writeJSON(w, http.StatusOK, pages[2])
		}
	})

	c := newTestClient(t, mux)

	var allIDs []int

	for page, err := range c.Blocklist.ListWithPagination(t.Context(), radarr.WithBlocklistPageSize(2)) {
		require.NoError(t, err)

		for _, r := range page.Records {
			allIDs = append(allIDs, r.ID)
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, allIDs)
}

func TestBlocklistService_ListByMovie(t *testing.T) {
	t.Parallel()

	want := []radarr.BlocklistRecord{
		{ID: 5, MovieID: 42, SourceTitle: new("Wanted.Movie.2024")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/blocklist/movie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "42", r.URL.Query().Get("movieId"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Blocklist.ListByMovie(t.Context(), 42)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 42, got[0].MovieID)
}

func TestBlocklistService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/blocklist/7", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Blocklist.Delete(t.Context(), 7)
	require.NoError(t, err)
}

func TestBlocklistService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/blocklist/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Blocklist.DeleteBulk(t.Context(), radarr.BlocklistBulkRequest{IDs: []int{1, 2, 3}})
	require.NoError(t, err)
}
