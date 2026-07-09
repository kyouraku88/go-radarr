package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQualityProfileService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.QualityProfile{
		{ID: 1, Name: new("HD-1080p"), UpgradeAllowed: true, Cutoff: 7},
		{ID: 2, Name: new("Ultra-HD"), UpgradeAllowed: false, Cutoff: 19},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/qualityprofile", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityProfile.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "HD-1080p", *got[0].Name)
	assert.True(t, got[0].UpgradeAllowed)
	assert.False(t, got[1].UpgradeAllowed)
}

func TestQualityProfileService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.QualityProfile{
		ID:             3,
		Name:           new("Any"),
		MinFormatScore: 10,
		Items: []radarr.QualityProfileItem{
			{ID: 1, Allowed: true},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/qualityprofile/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityProfile.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	assert.Equal(t, 10, got.MinFormatScore)
	require.Len(t, got.Items, 1)
	assert.True(t, got.Items[0].Allowed)
}

func TestQualityProfileService_Create(t *testing.T) {
	t.Parallel()

	input := radarr.QualityProfile{
		Name:           new("Custom"),
		UpgradeAllowed: true,
		Cutoff:         5,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/qualityprofile", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.QualityProfile
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotNil(t, body.Name)
		assert.Equal(t, "Custom", *body.Name)
		assert.True(t, body.UpgradeAllowed)
		body.ID = 99
		writeJSON(w, http.StatusCreated, body)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityProfile.Create(t.Context(), input)
	require.NoError(t, err)
	assert.Equal(t, 99, got.ID)
}

func TestQualityProfileService_Update(t *testing.T) {
	t.Parallel()

	input := radarr.QualityProfile{ID: 4, Name: new("Updated"), Cutoff: 8}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/qualityprofile/4", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.QualityProfile
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, 4, body.ID)
		writeJSON(w, http.StatusOK, body)
	})

	c := newTestClient(t, mux)
	got, err := c.QualityProfile.Update(t.Context(), 4, input)
	require.NoError(t, err)
	assert.Equal(t, 4, got.ID)
}

func TestQualityProfileService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/qualityprofile/6", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.QualityProfile.Delete(t.Context(), 6)
	require.NoError(t, err)
}
