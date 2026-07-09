package radarr_test

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueService_List(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.QueueRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 2,
		Records:      []radarr.QueueRecord{{ID: 1}, {ID: 2}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Queue.List(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, got.TotalRecords)
	require.Len(t, got.Records, 2)
	assert.Equal(t, 1, got.Records[0].ID)
}

func TestQueueService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "50", q.Get("pageSize"))
		assert.Equal(t, "date", q.Get("sortKey"))
		assert.Equal(t, "descending", q.Get("sortDirection"))
		assert.Equal(t, "true", q.Get("includeMovie"))
		assert.Equal(t, []string{"10", "11"}, q["movieIds"])
		assert.Equal(t, []string{"downloading", "completed"}, q["status"])
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.QueueRecord]{})
	})

	c := newTestClient(t, mux)
	_, err := c.Queue.List(t.Context(),
		radarr.WithQueuePage(2),
		radarr.WithQueuePageSize(50),
		radarr.WithQueueSortKey("date"),
		radarr.WithQueueSortDirection(radarr.SortDirectionDescending),
		radarr.WithQueueIncludeMovie(true),
		radarr.WithQueueMovieIDs(10, 11),
		radarr.WithQueueStatuses(radarr.QueueStatusDownloading, radarr.QueueStatusCompleted),
	)
	require.NoError(t, err)
}

func TestQueueService_ListWithPagination_MultiPage(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.QueueRecord]{
		{Page: 1, PageSize: 2, TotalRecords: 5, Records: []radarr.QueueRecord{{ID: 1}, {ID: 2}}},
		{Page: 2, PageSize: 2, TotalRecords: 5, Records: []radarr.QueueRecord{{ID: 3}, {ID: 4}}},
		{Page: 3, PageSize: 2, TotalRecords: 5, Records: []radarr.QueueRecord{{ID: 5}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue", func(w http.ResponseWriter, r *http.Request) {
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

	for page, err := range c.Queue.ListWithPagination(t.Context(), radarr.WithQueuePageSize(2)) {
		require.NoError(t, err)

		for _, r := range page.Records {
			allIDs = append(allIDs, r.ID)
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, allIDs)
}

func TestQueueService_ListWithPagination_ErrorOnPage(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue", func(w http.ResponseWriter, _ *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.QueueRecord]{
				Page: 1, PageSize: 2, TotalRecords: 4,
				Records: []radarr.QueueRecord{{ID: 1}, {ID: 2}},
			})
		} else {
			writeJSON(w, http.StatusInternalServerError, errorBody("db error"))
		}
	})

	c := newTestClient(t, mux)

	var gotErr error

	for _, err := range c.Queue.ListWithPagination(t.Context(), radarr.WithQueuePageSize(2)) {
		if err != nil {
			gotErr = err
			break
		}
	}

	require.Error(t, gotErr)

	var apiErr *radarr.APIError
	require.ErrorAs(t, gotErr, &apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
}

func TestQueueService_ListWithPagination_EarlyStop(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue", func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.QueueRecord]{
			Page: 1, PageSize: 2, TotalRecords: 100,
			Records: []radarr.QueueRecord{{ID: 1}, {ID: 2}},
		})
	})

	c := newTestClient(t, mux)

	for range c.Queue.ListWithPagination(t.Context(), radarr.WithQueuePageSize(2)) {
		break // stop after first page
	}

	// should only have called the server once
	assert.Equal(t, int32(1), calls.Load())
}

func TestQueueService_Status(t *testing.T) {
	t.Parallel()

	want := radarr.QueueStatusSummary{
		ID:         1,
		TotalCount: 10,
		Count:      8,
		Errors:     true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Queue.Status(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 10, got.TotalCount)
	assert.True(t, got.Errors)
}

func TestQueueService_ListDetails(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/queue/details", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "42", r.URL.Query().Get("movieId"))
		assert.Equal(t, "true", r.URL.Query().Get("includeMovie"))
		writeJSON(w, http.StatusOK, []radarr.QueueRecord{{ID: 5}})
	})

	c := newTestClient(t, mux)
	got, err := c.Queue.ListDetails(t.Context(),
		radarr.WithQueueDetailsMovieID(42),
		radarr.WithQueueDetailsIncludeMovie(true),
	)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 5, got[0].ID)
}

func TestQueueService_Grab(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/queue/grab/7", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Queue.Grab(t.Context(), 7)
	require.NoError(t, err)
}

func TestQueueService_GrabBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/queue/grab/bulk", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.QueueBulkRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{1, 2, 3}, body.IDs)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Queue.GrabBulk(t.Context(), radarr.QueueBulkRequest{IDs: []int{1, 2, 3}})
	require.NoError(t, err)
}

func TestQueueService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/queue/9", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "true", q.Get("removeFromClient"))
		assert.Equal(t, "true", q.Get("blocklist"))
		assert.Equal(t, "false", q.Get("skipRedownload"))
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Queue.Delete(t.Context(), 9,
		radarr.WithQueueRemoveFromClient(true),
		radarr.WithQueueBlocklist(true),
		radarr.WithQueueSkipRedownload(false),
	)
	require.NoError(t, err)
}

func TestQueueService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/queue/bulk", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.QueueBulkRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{4, 5, 6}, body.IDs)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Queue.DeleteBulk(t.Context(), radarr.QueueBulkRequest{IDs: []int{4, 5, 6}})
	require.NoError(t, err)
}
