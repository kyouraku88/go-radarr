package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryService_List(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []radarr.HistoryRecord{
			{ID: 1, MovieID: 42, EventType: radarr.MovieHistoryEventTypeGrabbed},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/history", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.History.List(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, got.TotalRecords)
	require.Len(t, got.Records, 1)
	assert.Equal(t, radarr.MovieHistoryEventTypeGrabbed, got.Records[0].EventType)
}

func TestHistoryService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/history", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "3", q.Get("page"))
		assert.Equal(t, "25", q.Get("pageSize"))
		assert.Equal(t, "date", q.Get("sortKey"))
		assert.Equal(t, "ascending", q.Get("sortDirection"))
		assert.Equal(t, "true", q.Get("includeMovie"))
		assert.Equal(t, "abc123", q.Get("downloadId"))
		assert.Equal(t, []string{"10", "20"}, q["movieIds"])
		assert.Equal(t, []string{"1", "3"}, q["eventType"])
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.HistoryRecord]{})
	})

	c := newTestClient(t, mux)
	_, err := c.History.List(t.Context(),
		radarr.WithHistoryPage(3),
		radarr.WithHistoryPageSize(25),
		radarr.WithHistorySortKey("date"),
		radarr.WithHistorySortDirection(radarr.SortDirectionAscending),
		radarr.WithHistoryIncludeMovie(true),
		radarr.WithHistoryDownloadID("abc123"),
		radarr.WithHistoryMovieIDs(10, 20),
		radarr.WithHistoryEventTypeIDs(1, 3),
	)
	require.NoError(t, err)
}

func TestHistoryService_ListWithPagination(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.HistoryRecord]{
				Page: 1, PageSize: 1, TotalRecords: 2,
				Records: []radarr.HistoryRecord{{ID: 1}},
			})
		case "2":
			writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.HistoryRecord]{
				Page: 2, PageSize: 1, TotalRecords: 2,
				Records: []radarr.HistoryRecord{{ID: 2}},
			})
		}
	})

	c := newTestClient(t, mux)

	var ids []int

	for page, err := range c.History.ListWithPagination(t.Context(), radarr.WithHistoryPageSize(1)) {
		require.NoError(t, err)

		for _, r := range page.Records {
			ids = append(ids, r.ID)
		}
	}

	assert.Equal(t, []int{1, 2}, ids)
}

func TestHistoryService_GetByMovie(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/history/movie", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "55", r.URL.Query().Get("movieId"))
		assert.Equal(t, string(radarr.MovieHistoryEventTypeDownloadFailed), r.URL.Query().Get("eventType"))
		assert.Equal(t, "true", r.URL.Query().Get("includeMovie"))
		writeJSON(w, http.StatusOK, []radarr.HistoryRecord{{ID: 9, MovieID: 55}})
	})

	c := newTestClient(t, mux)
	got, err := c.History.GetByMovie(t.Context(), 55,
		radarr.WithMovieHistoryEventType(radarr.MovieHistoryEventTypeDownloadFailed),
		radarr.WithMovieHistoryIncludeMovie(true),
	)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 9, got[0].ID)
	assert.Equal(t, 55, got[0].MovieID)
}

func TestHistoryService_Since(t *testing.T) {
	t.Parallel()

	since := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/history/since", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "2024-06-01T00:00:00Z", r.URL.Query().Get("date"))
		assert.Equal(t, string(radarr.MovieHistoryEventTypeGrabbed), r.URL.Query().Get("eventType"))
		writeJSON(w, http.StatusOK, []radarr.HistoryRecord{{ID: 3}})
	})

	c := newTestClient(t, mux)
	got, err := c.History.Since(t.Context(), since,
		radarr.WithSinceEventType(radarr.MovieHistoryEventTypeGrabbed),
	)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 3, got[0].ID)
}

func TestHistoryService_MarkFailed(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/history/failed/12", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.History.MarkFailed(t.Context(), 12)
	require.NoError(t, err)
}
