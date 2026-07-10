package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLanguageService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.LanguageResource{
		{ID: 1, Name: new("English"), NameLower: new("english")},
		{ID: 2, Name: new("German"), NameLower: new("german")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/language", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Language.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "English", *got[0].Name)
}

func TestLanguageService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.LanguageResource{ID: 3, Name: new("French")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/language/3", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Language.Get(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, 3, got.ID)
	require.NotNil(t, got.Name)
	assert.Equal(t, "French", *got.Name)
}
