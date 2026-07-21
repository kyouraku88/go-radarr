package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieService_List(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 1, TmdbID: 100}, {ID: 2, TmdbID: 200}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.List(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, want[0].ID, got[0].ID)
	assert.Equal(t, want[1].TmdbID, got[1].TmdbID)
}

func TestMovieService_List_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "12345", q.Get("tmdbId"))
		assert.Equal(t, "1", q.Get("languageId"))
		assert.Equal(t, "true", q.Get("excludeLocalCovers"))
		writeJSON(w, http.StatusOK, []radarr.Movie{})
	})

	c := newTestClient(t, mux)
	_, err := c.Movie.List(t.Context(),
		radarr.WithTmdbID(12345),
		radarr.WithLanguageID(1),
		radarr.WithExcludeLocalCovers(true),
	)
	require.NoError(t, err)
}

func TestMovieService_Get(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 42, TmdbID: 550, Title: new("Fight Club")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/42", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.Get(t.Context(), 42)
	require.NoError(t, err)
	assert.Equal(t, 42, got.ID)
	assert.Equal(t, 550, got.TmdbID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Fight Club", *got.Title)
}

func TestMovieService_Create(t *testing.T) {
	t.Parallel()

	input := radarr.Movie{TmdbID: 550, QualityProfileID: 1}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/movie", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.Movie
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, 550, body.TmdbID)
		assert.Equal(t, 1, body.QualityProfileID)
		body.ID = 10
		writeJSON(w, http.StatusCreated, body)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.Create(t.Context(), input)
	require.NoError(t, err)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, 550, got.TmdbID)
}

func TestMovieService_Update(t *testing.T) {
	t.Parallel()

	input := radarr.Movie{ID: 5, TmdbID: 550, Monitored: true}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/movie/5", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.Movie
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, 5, body.ID)
		assert.True(t, body.Monitored)
		writeJSON(w, http.StatusOK, body)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.Update(t.Context(), 5, input)
	require.NoError(t, err)
	assert.Equal(t, 5, got.ID)
}

func TestMovieService_Update_MoveFiles(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/movie/7", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("moveFiles"))
		writeJSON(w, http.StatusOK, radarr.Movie{ID: 7})
	})

	c := newTestClient(t, mux)
	_, err := c.Movie.Update(t.Context(), 7, radarr.Movie{ID: 7}, radarr.WithMoveFiles(true))
	require.NoError(t, err)
}

func TestMovieService_Delete(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/movie/3", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Movie.Delete(t.Context(), 3)
	require.NoError(t, err)
}

func TestMovieService_Delete_Options(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/movie/3", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "true", q.Get("deleteFiles"))
		assert.Equal(t, "true", q.Get("addImportExclusion"))
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.Movie.Delete(t.Context(), 3,
		radarr.WithDeleteFiles(true),
		radarr.WithAddImportExclusion(true),
	)
	require.NoError(t, err)
}

func TestMovieService_Lookup(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{
		{ID: 0, TmdbID: 550, Title: new("Fight Club")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/lookup", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "fight club", r.URL.Query().Get("term"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.Lookup(t.Context(), radarr.WithLookupTerm("fight club"))
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "Fight Club", *got[0].Title)
}

func TestMovieService_LookupByIMDB(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 0, TmdbID: 550, Title: new("Fight Club")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/lookup/imdb", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "tt0137523", r.URL.Query().Get("imdbId"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.LookupByIMDB(t.Context(), "tt0137523")
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Title)
	assert.Equal(t, "Fight Club", *got[0].Title)
}

func TestMovieService_LookupByTMDB(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 0, TmdbID: 550, Title: new("Fight Club")}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/lookup/tmdb", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		assert.Equal(t, "550", r.URL.Query().Get("tmdbId"))
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.LookupByTMDB(t.Context(), 550)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 550, got[0].TmdbID)
}

func TestMovieService_Import(t *testing.T) {
	t.Parallel()

	body := []radarr.Movie{
		{TmdbID: 550, QualityProfileID: 1, Title: new("Fight Club")},
		{TmdbID: 157336, QualityProfileID: 1, Title: new("Interstellar")},
	}
	want := []radarr.Movie{
		{ID: 10, TmdbID: 550, Title: new("Fight Club")},
		{ID: 11, TmdbID: 157336, Title: new("Interstellar")},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v3/movie/import", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.Import(t.Context(), body)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, 10, got[0].ID)
	assert.Equal(t, 11, got[1].ID)
}

func TestMovieService_GetFolder(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 42, TmdbID: 550, Title: new("Fight Club")}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/42/folder", func(w http.ResponseWriter, r *http.Request) {
		assertAPIKey(t, r)
		writeJSON(w, http.StatusOK, want)
	})

	c := newTestClient(t, mux)
	got, err := c.Movie.GetFolder(t.Context(), 42)
	require.NoError(t, err)
	assert.Equal(t, 42, got.ID)
	require.NotNil(t, got.Title)
	assert.Equal(t, "Fight Club", *got.Title)
}

func TestMovieService_Get_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		wantCode   int
	}{
		{"not found", http.StatusNotFound, http.StatusNotFound},
		{"server error", http.StatusInternalServerError, http.StatusInternalServerError},
		{"unauthorized", http.StatusUnauthorized, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/v3/movie/1", func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, tt.statusCode, errorBody("some error"))
			})

			c := newTestClient(t, mux)
			_, err := c.Movie.Get(t.Context(), 1)
			require.Error(t, err)

			var apiErr *radarr.APIError
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, tt.wantCode, apiErr.StatusCode)
		})
	}
}
