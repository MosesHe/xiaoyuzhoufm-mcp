package xyzclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"xiaoyuzhoufm-mcp/internal/constants"
)

// doSearch is a generic helper function to perform search requests.
// It returns the raw 'data' part of the response, highlight word, and load more key.
func doSearch(requestData SearchRequest) (json.RawMessage, *HighlightWord, *SearchAPILoadMoreKey, error) {
	apiURL := fmt.Sprintf("%s/v1/search/create", constants.APIBaseURL)
	slog.Debug("Performing search request", "url", apiURL, "type", requestData.Type, "keyword", requestData.Keyword)

	requestBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal request body for search: %w", err)
	}
	slog.Debug("Search request body", "body", string(requestBodyBytes))

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create request for search: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get token manager for search: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get access token for search: %w", err)
	}

	now := time.Now()
	isoTime := now.Format(time.RFC3339)

	req.Header.Set("Host", "api.xiaoyuzhoufm.com")
	req.Header.Set("User-Agent", "Xiaoyuzhou/2.57.1 (build:1576; iOS 17.4.1)")
	req.Header.Set("Market", "AppStore")
	req.Header.Set("App-BuildNo", "1576")
	req.Header.Set("OS", "ios")
	req.Header.Set("x-jike-access-token", accessToken)
	req.Header.Set("x-jike-device-id", constants.FixedDeviceID)
	req.Header.Set("Manufacturer", "Apple")
	req.Header.Set("BundleID", "app.podcast.cosmos")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("abtest-info", `{"old_user_discovery_feed":"enable"}`)
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1.0")
	req.Header.Set("Model", "iPhone14,2")
	req.Header.Set("app-permissions", "4")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Version", "2.57.1")
	req.Header.Set("WifiConnected", "true")
	req.Header.Set("OS-Version", "17.4.1")
	req.Header.Set("x-custom-xiaoyuzhou-app-dev", "")
	req.Header.Set("Local-Time", isoTime)
	req.Header.Set("Timezone", "Asia/Shanghai")

	slog.Debug("Sending HTTP request to search API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("http request failed for search: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read response body from search: %w", err)
	}
	slog.Debug("Received response from search API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("API request search failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	// Temporarily unmarshal into a structure that captures all top-level fields
	// and leaves 'data' as raw JSON for later specific parsing.
	var genericResponse struct {
		Data          json.RawMessage       `json:"data"`
		HighlightWord *HighlightWord        `json:"highlightWord,omitempty"`
		LoadMoreKey   *SearchAPILoadMoreKey `json:"loadMoreKey,omitempty"`
	}

	if err := json.Unmarshal(responseBodyBytes, &genericResponse); err != nil {
		slog.Error("Failed to unmarshal generic search response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, nil, nil, fmt.Errorf("failed to unmarshal generic search response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	return genericResponse.Data, genericResponse.HighlightWord, genericResponse.LoadMoreKey, nil
}

// SearchPodcasts searches for podcasts.
func SearchPodcasts(keyword string, loadMoreKey *SearchAPILoadMoreKey) (*PodcastSearchResponse, error) {
	request := SearchRequest{
		Keyword:     keyword,
		Type:        "PODCAST",
		LoadMoreKey: loadMoreKey,
	}
	rawData, highlight, lmk, err := doSearch(request)
	if err != nil {
		return nil, err
	}

	var podcasts []PodcastSearchResultItem
	if err := json.Unmarshal(rawData, &podcasts); err != nil {
		slog.Error("Failed to unmarshal podcast search data", "error", err, "rawData", string(rawData))
		return nil, fmt.Errorf("failed to unmarshal podcast search data: %w", err)
	}

	return &PodcastSearchResponse{
		Data:          podcasts,
		HighlightWord: highlight,
		LoadMoreKey:   lmk,
	}, nil
}

// SearchEpisodes searches for episodes.
func SearchEpisodes(keyword string, pid string, loadMoreKey *SearchAPILoadMoreKey) (*EpisodeSearchResponse, error) {
	request := SearchRequest{
		Keyword:     keyword,
		Type:        "EPISODE",
		PID:         pid,
		LoadMoreKey: loadMoreKey,
	}
	rawData, highlight, lmk, err := doSearch(request)
	if err != nil {
		return nil, err
	}

	var episodes []EpisodeSearchResultItem
	if err := json.Unmarshal(rawData, &episodes); err != nil {
		slog.Error("Failed to unmarshal episode search data", "error", err, "rawData", string(rawData))
		return nil, fmt.Errorf("failed to unmarshal episode search data: %w", err)
	}

	return &EpisodeSearchResponse{
		Data:          episodes,
		HighlightWord: highlight,
		LoadMoreKey:   lmk,
	}, nil
}

// SearchUsers searches for users.
func SearchUsers(keyword string, loadMoreKey *SearchAPILoadMoreKey) (*UserSearchResponse, error) {
	request := SearchRequest{
		Keyword:     keyword,
		Type:        "USER",
		LoadMoreKey: loadMoreKey,
	}
	rawData, highlight, lmk, err := doSearch(request)
	if err != nil {
		return nil, err
	}

	var users []UserSearchResultItem
	if err := json.Unmarshal(rawData, &users); err != nil {
		slog.Error("Failed to unmarshal user search data", "error", err, "rawData", string(rawData))
		return nil, fmt.Errorf("failed to unmarshal user search data: %w", err)
	}

	return &UserSearchResponse{
		Data:          users,
		HighlightWord: highlight,
		LoadMoreKey:   lmk,
	}, nil
}
