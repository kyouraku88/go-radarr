package radarr_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMovieEditorService_Edit(t *testing.T) {
	t.Parallel()

	input := radarr.MovieEditorRequest{
		MovieIDs:  []int{1, 2, 3},
		Monitored: new(true),
		Tags:      []int{10, 20},
		ApplyTags: radarr.ApplyTagsAdd,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/v3/movie/editor", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.MovieEditorRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{1, 2, 3}, body.MovieIDs)
		assert.NotNil(t, body.Monitored)
		assert.True(t, *body.Monitored)
		assert.Equal(t, radarr.ApplyTagsAdd, body.ApplyTags)
		w.WriteHeader(http.StatusAccepted)
	})

	c := newTestClient(t, mux)
	err := c.MovieEditor.Edit(t.Context(), input)
	require.NoError(t, err)
}

func TestMovieEditorService_Delete(t *testing.T) {
	t.Parallel()

	input := radarr.MovieEditorRequest{
		MovieIDs:    []int{4, 5},
		DeleteFiles: true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v3/movie/editor", func(w http.ResponseWriter, r *http.Request) {
		var body radarr.MovieEditorRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{4, 5}, body.MovieIDs)
		assert.True(t, body.DeleteFiles)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, mux)
	err := c.MovieEditor.Delete(t.Context(), input)
	require.NoError(t, err)
}
