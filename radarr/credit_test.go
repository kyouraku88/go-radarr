package radarr_test

import (
	"net/http"
	"testing"

	"github.com/bsido/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreditService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Credit{
		{ID: 1, PersonName: new("Brad Pitt"), Type: radarr.CreditTypeCast},
		{ID: 2, PersonName: new("David Fincher"), Type: radarr.CreditTypeCrew},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/credit", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Credit.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotNil(t, got[0].PersonName)
	assert.Equal(t, "Brad Pitt", *got[0].PersonName)
	assert.Equal(t, radarr.CreditTypeCast, got[0].Type)
}

func TestCreditService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/credit", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "55", q.Get("movieId"))
		assert.Equal(t, "3", q.Get("movieMetadataId"))
		writeJSON(w, http.StatusOK, []radarr.Credit{})
	})

	c := newTestClient(t, mux)
	_, err := c.Credit.List(t.Context(),
		radarr.WithCreditMovieID(55),
		radarr.WithCreditMovieMetadataID(3),
	)
	require.NoError(t, err)
}

func TestCreditService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Credit{ID: 9, PersonName: new("Edward Norton"), Type: radarr.CreditTypeCast}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/credit/9", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Credit.Get(t.Context(), 9)
	require.NoError(t, err)
	assert.Equal(t, 9, got.ID)
	assert.Equal(t, radarr.CreditTypeCast, got.Type)
}
