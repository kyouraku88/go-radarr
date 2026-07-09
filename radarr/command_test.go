package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.CommandRecord{
		{ID: 1, Name: new("RefreshMovie"), Status: radarr.CommandStatusCompleted},
		{ID: 2, Name: new("MissingMoviesSearch"), Status: radarr.CommandStatusQueued},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/command", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Command.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "RefreshMovie", *got[0].Name)
	assert.Equal(t, radarr.CommandStatusCompleted, got[0].Status)
	assert.Equal(t, radarr.CommandStatusQueued, got[1].Status)
}

func TestCommandService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.CommandRecord{
		ID:      5,
		Name:    new("RescanMovie"),
		Status:  radarr.CommandStatusStarted,
		Result:  radarr.CommandResultUnknown,
		Trigger: radarr.CommandTriggerManual,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/command/5", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Command.Get(t.Context(), 5)
	require.NoError(t, err)
	assert.Equal(t, 5, got.ID)
	assert.Equal(t, radarr.CommandStatusStarted, got.Status)
	assert.Equal(t, radarr.CommandTriggerManual, got.Trigger)
}

func TestCommandService_Create(t *testing.T) {
	t.Parallel()

	input := radarr.CommandRecord{
		Name:    new("RefreshMovie"),
		Trigger: radarr.CommandTriggerManual,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/command", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.CommandRecord
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotNil(t, body.Name)
		assert.Equal(t, "RefreshMovie", *body.Name)
		body.ID = 10
		body.Status = radarr.CommandStatusQueued
		writeJSON(w, http.StatusCreated, body)
	})

	c := newTestClient(t, mux)
	got, err := c.Command.Create(t.Context(), input)
	require.NoError(t, err)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, radarr.CommandStatusQueued, got.Status)
}

func TestCommandService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/command/3", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Command.Delete(t.Context(), 3)
	require.NoError(t, err)
}

func TestCommandService_Create_Error(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/command", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusBadRequest, errorBody("Unknown command name"))
	})

	c := newTestClient(t, mux)
	_, err := c.Command.Create(t.Context(), radarr.CommandRecord{Name: new("BadCommand")})
	require.Error(t, err)

	var apiErr *radarr.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "Unknown command name", apiErr.Message)
}
