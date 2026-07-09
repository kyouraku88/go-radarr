# go-radarr

A Go client library for the [Radarr](https://radarr.video/docs/api/) v3 API.

## Installation

```sh
go get github.com/kyouraku88/go-radarr
```

## Usage

```go
client, err := radarr.New(
    radarr.WithBaseURL("http://localhost:7878"),
    radarr.WithAPIKey("your-api-key"),
)
if err != nil {
    log.Fatal(err)
}

// Fetch a single movie
movie, err := client.Movie.Get(ctx, 42)

// List all movies
movies, err := client.Movie.List(ctx, radarr.WithTmdbID(550))

// Add a movie
created, err := client.Movie.Create(ctx, radarr.Movie{
    TmdbID:           550,
    QualityProfileID: 1,
    Monitored:        true,
})

// Delete with options
err = client.Movie.Delete(ctx, 42, radarr.WithDeleteFiles(true))
```

### Pagination

Paginated endpoints (Queue, History) expose two methods:

```go
// Single page — explicit control
page, err := client.Queue.List(ctx, radarr.WithPage(2), radarr.WithPageSize(50))

// All pages — automatic iteration via range-over-func (Go 1.23+)
for page, err := range client.Queue.ListWithPagination(ctx, radarr.WithPageSize(100)) {
    if err != nil {
        break
    }
    for _, record := range page.Records { ... }
}
```

### Custom HTTP client

```go
client, err := radarr.New(
    radarr.WithBaseURL("http://localhost:7878"),
    radarr.WithAPIKey("your-api-key"),
    radarr.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

### Error handling

HTTP errors are returned as `*radarr.APIError`, which carries the status code:

```go
var apiErr *radarr.APIError
if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
    // handle 404
}
```

## Dependencies

- [`github.com/go-resty/resty/v2`](https://github.com/go-resty/resty) — HTTP client
