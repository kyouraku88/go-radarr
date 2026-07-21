package radarr_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyouraku88/go-radarr/radarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// newTestClient starts an httptest.Server backed by mux and returns a Radarr
// client pointed at it. The server is closed automatically when t ends.
func newTestClient(t *testing.T, mux *http.ServeMux) *radarr.Radarr {
	t.Helper()

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	c, err := radarr.New(radarr.WithBaseURL(srv.URL), radarr.WithAPIKey("test-key"))
	require.NoError(t, err)

	return c
}

// writeJSON encodes v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// errorBody returns a Radarr-style error response body.
func errorBody(msg string) []map[string]string {
	return []map[string]string{{"errorMessage": msg}}
}

// assertAPIKey verifies the X-Api-Key header is present in the request.
func assertAPIKey(t *testing.T, r *http.Request) {
	t.Helper()
	assert.Equal(t, "test-key", r.Header.Get("X-Api-Key"))
}

// ---------------------------------------------------------------------------
// TestNew — client construction validation
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []radarr.ClientOption
		wantErr string
	}{
		{
			name:    "missing base URL",
			opts:    []radarr.ClientOption{radarr.WithAPIKey("key")},
			wantErr: "base URL is required",
		},
		{
			name:    "missing API key",
			opts:    []radarr.ClientOption{radarr.WithBaseURL("http://localhost")},
			wantErr: "API key is required",
		},
		{
			name: "valid",
			opts: []radarr.ClientOption{
				radarr.WithBaseURL("http://localhost:7878"),
				radarr.WithAPIKey("abc123"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := radarr.New(tt.opts...)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, c)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, c)
		})
	}
}

// ---------------------------------------------------------------------------
// TestAPIError — typed error returned on non-2xx responses
// ---------------------------------------------------------------------------

func TestAPIError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  *radarr.APIError
		want string
	}{
		{
			name: "with message",
			err:  &radarr.APIError{StatusCode: 404, Status: "404 Not Found", Message: "Movie not found"},
			want: "radarr: HTTP 404: Movie not found",
		},
		{
			name: "without message",
			err:  &radarr.APIError{StatusCode: 500, Status: "500 Internal Server Error"},
			want: "radarr: HTTP 500 Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestAPIError_NonJSON_Body(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/1", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	})

	c := newTestClient(t, mux)
	_, err := c.Movie.Get(t.Context(), 1)

	require.Error(t, err)

	var apiErr *radarr.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	assert.Empty(t, apiErr.Message) // non-JSON body → no message extracted
}

func TestAPIError_JSONBody(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/movie/99", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusNotFound, errorBody("Movie with id 99 does not exist"))
	})

	c := newTestClient(t, mux)
	_, err := c.Movie.Get(t.Context(), 99)

	require.Error(t, err)

	var apiErr *radarr.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "Movie with id 99 does not exist", apiErr.Message)
}
