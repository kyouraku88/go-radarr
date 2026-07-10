package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexerFlagService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.IndexerFlag{
		{ID: 1, Name: new("G_Freeleech"), NameLower: new("g_freeleech")},
		{ID: 2, Name: new("G_Halfleech"), NameLower: new("g_halfleech")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/indexerflag", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.IndexerFlag.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "G_Freeleech", *got[0].Name)
}
