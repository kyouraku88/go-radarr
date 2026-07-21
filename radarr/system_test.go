package radarr_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemService_Status(t *testing.T) {
	t.Parallel()

	want := radarr.SystemStatus{
		AppName:      new("Radarr"),
		Version:      new("5.0.0"),
		IsDocker:     true,
		DatabaseType: radarr.DatabaseTypeSQLite,
		BuildTime:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		StartTime:    time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/system/status", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.Status(t.Context())
	require.NoError(t, err)
	require.NotNil(t, got.AppName)
	assert.Equal(t, "Radarr", *got.AppName)
	require.NotNil(t, got.Version)
	assert.Equal(t, "5.0.0", *got.Version)
	assert.True(t, got.IsDocker)
	assert.Equal(t, radarr.DatabaseTypeSQLite, got.DatabaseType)
}

func TestSystemService_ListTasks(t *testing.T) {
	t.Parallel()

	want := []radarr.Task{
		{ID: 1, Name: new("RefreshMovie"), Interval: 360},
		{ID: 2, Name: new("BackupDatabase"), Interval: 1440},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/system/task", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.ListTasks(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, 1, got[0].ID)
	require.NotNil(t, got[1].Name)
	assert.Equal(t, "BackupDatabase", *got[1].Name)
}

func TestSystemService_GetTask(t *testing.T) {
	t.Parallel()

	want := radarr.Task{ID: 3, Name: new("MessagingCleanup"), Interval: 60}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/system/task/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.GetTask(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.Equal(t, 60, got.Interval)
}

func TestSystemService_ListBackups(t *testing.T) {
	t.Parallel()

	want := []radarr.Backup{
		{ID: 1, Name: new("radarr.db"), Type: radarr.BackupTypeScheduled, Size: 1024},
		{ID: 2, Name: new("radarr_manual.db"), Type: radarr.BackupTypeManual, Size: 2048},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/system/backup", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.ListBackups(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, radarr.BackupTypeScheduled, got[0].Type)
	assert.Equal(t, int64(2048), got[1].Size)
}

func TestSystemService_DeleteBackup(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/system/backup/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.System.DeleteBackup(t.Context(), 5)
	require.NoError(t, err)
}

func TestSystemService_RestoreBackup(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/system/backup/restore/2", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.System.RestoreBackup(t.Context(), 2)
	require.NoError(t, err)
}

func TestSystemService_Restart(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/system/restart", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.System.Restart(t.Context())
	require.NoError(t, err)
}

func TestSystemService_Shutdown(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/system/shutdown", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.System.Shutdown(t.Context())
	require.NoError(t, err)
}

func TestSystemService_Health(t *testing.T) {
	t.Parallel()

	want := []radarr.HealthRecord{
		{ID: 1, Source: new("IndexerStatusCheck"), Type: radarr.HealthCheckResultWarning, Message: new("Indexer unavailable")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.Health(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, radarr.HealthCheckResultWarning, got[0].Type)
	require.NotNil(t, got[0].Message)
	assert.Equal(t, "Indexer unavailable", *got[0].Message)
}

func TestSystemService_DiskSpace(t *testing.T) {
	t.Parallel()

	want := []radarr.DiskSpace{
		{ID: 1, Path: new("/mnt/media"), FreeSpace: 500_000_000_000, TotalSpace: 2_000_000_000_000},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/diskspace", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.DiskSpace(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, int64(500_000_000_000), got[0].FreeSpace)
}

func TestSystemService_Updates(t *testing.T) {
	t.Parallel()

	want := []radarr.Update{
		{
			ID:          1,
			Version:     new("5.1.0"),
			Installed:   false,
			Installable: true,
			Latest:      true,
			Changes:     &radarr.UpdateChanges{New: []string{"Feature A"}, Fixed: []string{"Bug B"}},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/update", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.System.Updates(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].Latest)
	require.NotNil(t, got[0].Changes)
	assert.Equal(t, []string{"Feature A"}, got[0].Changes.New)
	assert.Equal(t, []string{"Bug B"}, got[0].Changes.Fixed)
}
