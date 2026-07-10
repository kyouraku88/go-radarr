package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExclusionsService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.ImportListExclusion{
		{ID: 1, TmdbID: 111, MovieTitle: new("Batman v Superman"), MovieYear: 2016},
		{ID: 2, TmdbID: 222, MovieTitle: new("Justice League"), MovieYear: 2017},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/exclusions", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Exclusions.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].MovieTitle)
	assert.Equal(t, "Batman v Superman", *got[0].MovieTitle)
}

func TestExclusionsService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.ImportListExclusion{ID: 3, TmdbID: 333, MovieTitle: new("Aquaman"), MovieYear: 2018}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/exclusions/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Exclusions.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.Equal(t, 2018, got.MovieYear)
}

func TestExclusionsService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.ImportListExclusion{TmdbID: 444, MovieTitle: new("Shazam"), MovieYear: 2019}
	want := radarr.ImportListExclusion{ID: 4, TmdbID: 444, MovieTitle: new("Shazam"), MovieYear: 2019}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/exclusions", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Exclusions.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestExclusionsService_CreateBulk(t *testing.T) {
	t.Parallel()

	body := []radarr.ImportListExclusion{
		{TmdbID: 555, MovieTitle: new("Birds of Prey"), MovieYear: 2020},
		{TmdbID: 666, MovieTitle: new("Wonder Woman 1984"), MovieYear: 2020},
	}
	want := []radarr.ImportListExclusion{
		{ID: 5, TmdbID: 555, MovieTitle: new("Birds of Prey"), MovieYear: 2020},
		{ID: 6, TmdbID: 666, MovieTitle: new("Wonder Woman 1984"), MovieYear: 2020},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/exclusions/bulk", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Exclusions.CreateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, 5, got[0].ID)
}

func TestExclusionsService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/exclusions/7", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Exclusions.Delete(t.Context(), 7)
	require.NoError(t, err)
}

func TestExclusionsService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/exclusions/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Exclusions.DeleteBulk(t.Context(), radarr.ExclusionsBulkRequest{IDs: []int{1, 2, 3}})
	require.NoError(t, err)
}

func TestExclusionsService_ListPaged(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.ImportListExclusion]{
		Page: 1, PageSize: 10, TotalRecords: 1,
		Records: []radarr.ImportListExclusion{{ID: 1, TmdbID: 111}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/exclusions/paged", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Exclusions.ListPaged(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.TotalRecords)
	require.Len(t, got.Records, 1)
	assert.Equal(t, 111, got.Records[0].TmdbID)
}

func TestExclusionsService_ListPagedWithPagination(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.ImportListExclusion]{
		{Page: 1, PageSize: 2, TotalRecords: 4, Records: []radarr.ImportListExclusion{{ID: 1, TmdbID: 111}, {ID: 2, TmdbID: 222}}},
		{Page: 2, PageSize: 2, TotalRecords: 4, Records: []radarr.ImportListExclusion{{ID: 3, TmdbID: 333}, {ID: 4, TmdbID: 444}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/exclusions/paged", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, pages[0])
		default:
			writeJSON(w, http.StatusOK, pages[1])
		}
	})

	c := newTestClient(t, mux)

	var total int

	for page, err := range c.Exclusions.ListPagedWithPagination(t.Context(), radarr.WithExclusionsPageSize(2)) {
		require.NoError(t, err)

		total += len(page.Records)
	}

	assert.Equal(t, 4, total)
}
