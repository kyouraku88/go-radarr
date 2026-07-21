package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomFormatService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.CustomFormat{
		{ID: 1, Name: new("Remux"), IncludeCustomFormatWhenRenaming: new(true)},
		{ID: 2, Name: new("HDR"), IncludeCustomFormatWhenRenaming: new(false)},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/customformat", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "Remux", *got[0].Name)
}

func TestCustomFormatService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.CustomFormat{ID: 3, Name: new("Dolby Vision")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/customformat/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Name)
	assert.Equal(t, "Dolby Vision", *got.Name)
}

func TestCustomFormatService_Schema(t *testing.T) {
	t.Parallel()

	want := []radarr.CustomFormat{{ID: 0, Name: new("ReleaseTitleSpecification")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/customformat/schema", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.Schema(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestCustomFormatService_Create(t *testing.T) {
	t.Parallel()

	body := radarr.CustomFormat{Name: new("Atmos")}
	want := radarr.CustomFormat{ID: 4, Name: new("Atmos")}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/customformat", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusCreated, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.Create(t.Context(), body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestCustomFormatService_Update(t *testing.T) {
	t.Parallel()

	body := radarr.CustomFormat{ID: 4, Name: new("Atmos Updated")}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/customformat/4", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.Update(t.Context(), 4, body)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestCustomFormatService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := radarr.CustomFormatBulkRequest{IDs: []int{1, 2}, IncludeCustomFormatWhenRenaming: new(true)}
	want := []radarr.CustomFormat{
		{ID: 1, Name: new("Remux")},
		{ID: 2, Name: new("HDR")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/customformat/bulk", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.CustomFormat.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestCustomFormatService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/customformat/5", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.CustomFormat.Delete(t.Context(), 5)
	require.NoError(t, err)
}

func TestCustomFormatService_DeleteBulk(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/customformat/bulk", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.CustomFormat.DeleteBulk(t.Context(), radarr.CustomFormatBulkRequest{IDs: []int{1, 2}})
	require.NoError(t, err)
}
