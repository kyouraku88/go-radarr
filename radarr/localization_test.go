package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalizationService_Get(t *testing.T) {
	t.Parallel()

	want := map[string]string{
		"Add":    "Add",
		"Cancel": "Cancel",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/localization", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(want)
	})

	c := newTestClient(t, mux)
	got, err := c.Localization.Get(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "Add", got["Add"])
	assert.Equal(t, "Cancel", got["Cancel"])
}

func TestLocalizationService_GetLanguage(t *testing.T) {
	t.Parallel()

	want := radarr.LocalizationLanguage{Identifier: new("en")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/localization/language", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Localization.GetLanguage(t.Context())
	require.NoError(t, err)
	require.NotNil(t, got.Identifier)
	assert.Equal(t, "en", *got.Identifier)
}
