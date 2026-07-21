package radarr_test

import (
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQualityDefinitionService_List(t *testing.T) {
	t.Parallel()

	min1 := 0.0
	max1 := 100.0
	want := []radarr.QualityDefinition{
		{ID: 1, Title: new("SDTV"), Weight: 1, MinSize: &min1, MaxSize: &max1},
		{ID: 2, Title: new("720p"), Weight: 2},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/qualitydefinition", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityDefinition.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "SDTV", *got[0].Title)
	assert.Equal(t, 1, got[0].Weight)
}

func TestQualityDefinitionService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.QualityDefinition{ID: 3, Title: new("1080p"), Weight: 3}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/qualitydefinition/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityDefinition.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "1080p", *got.Title)
}

func TestQualityDefinitionService_Update(t *testing.T) {
	t.Parallel()

	maxSize := 200.0
	body := radarr.QualityDefinition{ID: 3, Title: new("1080p"), Weight: 3, MaxSize: &maxSize}
	want := body

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/qualitydefinition/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityDefinition.Update(t.Context(), 3, body)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.MaxSize)
	assert.InDelta(t, 200.0, *got.MaxSize, 0.001)
}

func TestQualityDefinitionService_UpdateBulk(t *testing.T) {
	t.Parallel()

	body := []radarr.QualityDefinition{
		{ID: 1, Title: new("SDTV"), Weight: 1},
		{ID: 2, Title: new("720p"), Weight: 2},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/qualitydefinition/update", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, body)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityDefinition.UpdateBulk(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestQualityDefinitionService_Limits(t *testing.T) {
	t.Parallel()

	want := radarr.QualityDefinitionLimits{Min: 0, Max: 1200}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/qualitydefinition/limits", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityDefinition.Limits(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 0, got.Min)
	assert.Equal(t, 1200, got.Max)
}
