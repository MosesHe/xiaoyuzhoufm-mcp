package xyzclient

import (
	"bytes" // Added bytes import
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"xiaoyuzhoufm-mcp/internal/constants"
)

// GetPodcastDetailsByID fetches detailed information for a specific podcast by its PID.
func GetPodcastDetailsByID(podcastID string) (*PodcastDetailData, error) {
	if podcastID == "" {
		return nil, fmt.Errorf("podcastID cannot be empty")
	}
	apiURL := fmt.Sprintf("%s/v1/podcast/get?pid=%s", constants.APIBaseURL, podcastID)
	slog.Debug("Fetching podcast details by ID", "url", apiURL, "podcastID", podcastID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetPodcastDetailsByID: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get token manager: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set Headers - consistent with GetUserProfileByID
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
	req.Header.Set("Accept", "*/*")
	req.Header.Set("App-Version", "2.57.1")
	req.Header.Set("WifiConnected", "true")
	req.Header.Set("OS-Version", "17.4.1")
	req.Header.Set("x-custom-xiaoyuzhou-app-dev", "")
	req.Header.Set("Local-Time", isoTime)
	req.Header.Set("Timezone", "Asia/Shanghai")

	slog.Debug("Sending HTTP request to GetPodcastDetailsByID API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for GetPodcastDetailsByID: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from GetPodcastDetailsByID: %w", err)
	}
	slog.Debug("Received response from GetPodcastDetailsByID API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request GetPodcastDetailsByID failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	var responseWrapper PodcastDetailAPIResponse
	if err := json.Unmarshal(responseBodyBytes, &responseWrapper); err != nil {
		slog.Error("Failed to unmarshal GetPodcastDetailsByID success response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, fmt.Errorf("failed to unmarshal GetPodcastDetailsByID success response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	slog.Debug("Successfully fetched and parsed podcast details.", "podcastID", podcastID, "title", responseWrapper.Data.Title)
	return &responseWrapper.Data, nil
}

// ListPodcastEpisodes fetches a list of episodes for a specific podcast.
func ListPodcastEpisodes(requestData EpisodeListRequest) (*EpisodeListResponseData, error) {
	if requestData.PID == "" {
		return nil, fmt.Errorf("podcastID (PID) in requestData cannot be empty")
	}

	apiURL := fmt.Sprintf("%s/v1/episode/list", constants.APIBaseURL)
	slog.Debug("Fetching podcast episodes list", "url", apiURL, "podcastID", requestData.PID)

	requestBodyBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for ListPodcastEpisodes: %w", err)
	}
	slog.Debug("ListPodcastEpisodes request body", "body", string(requestBodyBytes))

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for ListPodcastEpisodes: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get token manager: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set Headers - consistent with other API calls
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
	req.Header.Set("Accept", "application/json")       // Specify JSON acceptance
	req.Header.Set("Content-Type", "application/json") // Specify JSON content type
	req.Header.Set("App-Version", "2.57.1")
	req.Header.Set("WifiConnected", "true")
	req.Header.Set("OS-Version", "17.4.1")
	req.Header.Set("x-custom-xiaoyuzhou-app-dev", "")
	req.Header.Set("Local-Time", isoTime)
	req.Header.Set("Timezone", "Asia/Shanghai")

	slog.Debug("Sending HTTP request to ListPodcastEpisodes API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for ListPodcastEpisodes: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from ListPodcastEpisodes: %w", err)
	}
	slog.Debug("Received response from ListPodcastEpisodes API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request ListPodcastEpisodes failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	var responseData EpisodeListResponseData                                 // Changed variable name and type
	if err := json.Unmarshal(responseBodyBytes, &responseData); err != nil { // Changed target of Unmarshal
		slog.Error("Failed to unmarshal ListPodcastEpisodes success response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, fmt.Errorf("failed to unmarshal ListPodcastEpisodes success response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	slog.Debug("Successfully fetched and parsed podcast episodes list.", "podcastID", requestData.PID, "count", len(responseData.Data)) // Changed to responseData.Data
	return &responseData, nil                                                                                                           // Changed to return &responseData
}

// GetEpisodeDetailsByID fetches detailed information for a specific episode by its EID.
func GetEpisodeDetailsByID(episodeID string) (*Episode, error) {
	if episodeID == "" {
		return nil, fmt.Errorf("episodeID cannot be empty")
	}
	apiURL := fmt.Sprintf("%s/v1/episode/get?eid=%s", constants.APIBaseURL, episodeID)
	slog.Debug("Fetching episode details by ID", "url", apiURL, "episodeID", episodeID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetEpisodeDetailsByID: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get token manager: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Set Headers - consistent with other API calls
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
	req.Header.Set("Accept", "application/json") // Prefer JSON response
	req.Header.Set("App-Version", "2.57.1")
	req.Header.Set("WifiConnected", "true")
	req.Header.Set("OS-Version", "17.4.1")
	req.Header.Set("x-custom-xiaoyuzhou-app-dev", "")
	req.Header.Set("Local-Time", isoTime)
	req.Header.Set("Timezone", "Asia/Shanghai")

	slog.Debug("Sending HTTP request to GetEpisodeDetailsByID API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for GetEpisodeDetailsByID: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from GetEpisodeDetailsByID: %w", err)
	}
	slog.Debug("Received response from GetEpisodeDetailsByID API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request GetEpisodeDetailsByID failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	var responseWrapper EpisodeDetailAPIResponse // Use the new wrapper
	if err := json.Unmarshal(responseBodyBytes, &responseWrapper); err != nil {
		slog.Error("Failed to unmarshal GetEpisodeDetailsByID success response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, fmt.Errorf("failed to unmarshal GetEpisodeDetailsByID success response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	slog.Debug("Successfully fetched and parsed episode details.", "episodeID", episodeID, "title", responseWrapper.Data.Title)
	return &responseWrapper.Data, nil
}
