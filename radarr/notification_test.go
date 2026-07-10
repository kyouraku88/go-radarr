package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Notification{
		{ID: 1, Name: new("Slack"), OnGrab: true, OnDownload: true},
		{ID: 2, Name: new("Discord"), OnGrab: false, OnDownload: true},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/notification", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Notification.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Slack", *got[0].Name)
	assert.True(t, got[0].OnGrab)
}

func TestNotificationService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Notification{ID: 3, Name: new("Pushover"), OnGrab: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/notification/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Notification.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
}

func TestNotificationService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.Notification{{ID: 0, Name: new("Slack")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/notification/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Notification.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestNotificationService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.Notification{Name: new("Email"), OnGrab: true}
	want := radarr.Notification{ID: 4, Name: new("Email"), OnGrab: true}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/notification", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Notification.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestNotificationService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.Notification{ID: 4, Name: new("Email"), OnGrab: false}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/notification/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Notification.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
	assert.False(t, got.OnGrab)
}

func TestNotificationService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/notification/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Notification.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestNotificationService_Test(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/notification/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Notification.Test(t.Context(), radarr.Notification{Name: new("Slack")})
	require.NoError(t, err)
}

func TestNotificationService_TestAll(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/notification/testall", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Notification.TestAll(t.Context())
	require.NoError(t, err)
}

func TestNotificationService_Action(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/notification/action/sendTest", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Notification.Action(t.Context(), "sendTest", radarr.Notification{Name: new("Slack")})
	require.NoError(t, err)
}
