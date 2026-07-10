package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogService_List(t *testing.T) {
	t.Parallel()

	want := radarr.PagedResult[radarr.LogRecord]{
		Page: 1, PageSize: 10, TotalRecords: 2,
		Records: []radarr.LogRecord{
			{ID: 1, Level: new("info"), Message: new("Task completed"), Time: time.Now()},
			{ID: 2, Level: new("warn"), Message: new("Disk low"), Time: time.Now()},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Log.List(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 2, got.TotalRecords)
	require.Len(t, got.Records, 2)
	require.NotNil(t, got.Records[0].Level)
	assert.Equal(t, "info", *got.Records[0].Level)
}

func TestLogService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "1", q.Get("page"))
		assert.Equal(t, "50", q.Get("pageSize"))
		assert.Equal(t, "time", q.Get("sortKey"))
		assert.Equal(t, "error", q.Get("level"))
		writeJSON(w, http.StatusOK, radarr.PagedResult[radarr.LogRecord]{})
	})

	c := newTestClient(t, mux)
	_, err := c.Log.List(t.Context(),
		radarr.WithLogPage(1),
		radarr.WithLogPageSize(50),
		radarr.WithLogSortKey("time"),
		radarr.WithLogLevel("error"),
	)
	require.NoError(t, err)
}

func TestLogService_ListWithPagination(t *testing.T) {
	t.Parallel()

	pages := []radarr.PagedResult[radarr.LogRecord]{
		{Page: 1, PageSize: 2, TotalRecords: 4, Records: []radarr.LogRecord{{ID: 1}, {ID: 2}}},
		{Page: 2, PageSize: 2, TotalRecords: 4, Records: []radarr.LogRecord{{ID: 3}, {ID: 4}}},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			writeJSON(w, http.StatusOK, pages[0])
		default:
			writeJSON(w, http.StatusOK, pages[1])
		}
	})

	c := newTestClient(t, mux)

	var allIDs []int

	for page, err := range c.Log.ListWithPagination(t.Context(), radarr.WithLogPageSize(2)) {
		require.NoError(t, err)

		for _, r := range page.Records {
			allIDs = append(allIDs, r.ID)
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4}, allIDs)
}

func TestLogService_ListFiles(t *testing.T) {
	t.Parallel()

	want := []radarr.LogFile{
		{ID: 1, Filename: new("radarr.txt"), LastWriteTime: time.Now()},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log/file", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Log.ListFiles(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Filename)
	assert.Equal(t, "radarr.txt", *got[0].Filename)
}

func TestLogService_GetFile(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log/file/radarr.txt", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("log content here"))
	})

	c := newTestClient(t, mux)
	got, err := c.Log.GetFile(t.Context(), "radarr.txt")
	require.NoError(t, err)
	assert.Equal(t, "log content here", got)
}

func TestLogService_ListUpdateFiles(t *testing.T) {
	t.Parallel()

	want := []radarr.LogFile{
		{ID: 1, Filename: new("update.txt")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log/file/update", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Log.ListUpdateFiles(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestLogService_GetUpdateFile(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/log/file/update/update.txt", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("update log content"))
	})

	c := newTestClient(t, mux)
	got, err := c.Log.GetUpdateFile(t.Context(), "update.txt")
	require.NoError(t, err)
	assert.Equal(t, "update log content", got)
}
