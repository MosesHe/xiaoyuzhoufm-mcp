package xyzclient

import (
	"net/http"
	"sync"
	"time"
)

var (
	httpClient *http.Client
	once       sync.Once
)

// GetHTTPClient returns a singleton HTTP client.
// The client is initialized once.
func GetHTTPClient() *http.Client {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: 30 * time.Second, // Default timeout
		}
	})
	return httpClient
}
