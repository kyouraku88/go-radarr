package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWantedService_ListCutoff(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.Movie]{
		Page: 1, PageSize: 10, TotalRecords: 2,
		Records: []radarr.Movie{
			{ID: 1, Title: new("Interstellar"), TmdbID: 157336},
			{ID: 2, Title: new("The Martian"), TmdbID: 286217},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/wanted/cutoff", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Wanted.ListCutoff(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, got.TotalRecords)
	require.Len(t, got.Records, 2)
	require.NotNil(t, got.Records[0].Title)
	assert.Equal(t, "Interstellar", *got.Records[0].Title)
}

func TestWantedService_ListCutoff_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/wanted/cutoff", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "20", q.Get("pageSize"))
		assert.Equal(t, "title", q.Get("sortKey"))
		assert.Equal(t, "true", q.Get("monitored"))
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.Movie]{})
	})

	c := newTestClient(t, mux)
	_, err := c.Wanted.ListCutoff(t.Context(),
		radarr.WithWantedCutoffPage(2),
		radarr.WithWantedCutoffPageSize(20),
		radarr.WithWantedCutoffSortKey("title"),
		radarr.WithWantedCutoffMonitored(true),
	)
	require.NoError(t, err)
}

func TestWantedService_ListCutoffWithPagination(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.Movie]{
		{Page: 1, PageSize: 2, TotalRecords: 4, Records: []radarr.Movie{{ID: 1, Title: new("Interstellar")}, {ID: 2, Title: new("The Martian")}}},
		{Page: 2, PageSize: 2, TotalRecords: 4, Records: []radarr.Movie{{ID: 3, Title: new("Dune")}, {ID: 4, Title: new("Arrival")}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/wanted/cutoff", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, pages[0])
		default:
			writeJSON(w, http.StatusOK, pages[1])
		}
	})

	c := newTestClient(t, mux)

	var allIDs []int

	for page, err := range c.Wanted.ListCutoffWithPagination(t.Context(), radarr.WithWantedCutoffPageSize(2)) {
		require.NoError(t, err)

		for _, m := range page.Records {
			allIDs = append(allIDs, m.ID)
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4}, allIDs)
}

func TestWantedService_ListMissing(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.Movie]{
		Page: 1, PageSize: 10, TotalRecords: 1,
		Records: []radarr.Movie{
			{ID: 3, Title: new("Dune"), TmdbID: 438631},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/wanted/missing", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Wanted.ListMissing(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.TotalRecords)
	require.Len(t, got.Records, 1)
	require.NotNil(t, got.Records[0].Title)
	assert.Equal(t, "Dune", *got.Records[0].Title)
}

func TestWantedService_ListMissingWithPagination(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.Movie]{
		{Page: 1, PageSize: 2, TotalRecords: 3, Records: []radarr.Movie{{ID: 3, Title: new("Dune")}, {ID: 4, Title: new("Dune: Part Two")}}},
		{Page: 2, PageSize: 2, TotalRecords: 3, Records: []radarr.Movie{{ID: 5, Title: new("Oppenheimer")}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/wanted/missing", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, pages[0])
		default:
			writeJSON(w, http.StatusOK, pages[1])
		}
	})

	c := newTestClient(t, mux)

	var titles []string

	for page, err := range c.Wanted.ListMissingWithPagination(t.Context(), radarr.WithWantedMissingPageSize(2)) {
		require.NoError(t, err)

		for _, m := range page.Records {
			require.NotNil(t, m.Title)

			titles = append(titles, *m.Title)
		}
	}

	assert.Equal(t, []string{"Dune", "Dune: Part Two", "Oppenheimer"}, titles)
}
