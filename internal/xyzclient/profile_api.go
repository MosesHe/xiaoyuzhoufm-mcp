package xyzclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"xiaoyuzhoufm-mcp/internal/constants"
)

// GetUserProfileByID fetches a user's public profile information by their UID.
func GetUserProfileByID(userID string) (*UserProfileData, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	apiURL := fmt.Sprintf("%s/v1/profile/get?uid=%s", constants.APIBaseURL, userID)
	slog.Debug("Fetching user profile by ID", "url", apiURL, "userID", userID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetUserProfileByID: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get token manager: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	// Set Headers - based on temp/xyz/handlers/profile.go -> GetProfileByUid
	now := time.Now()
	// Matches "yyyy-MM-dd'T'HH:mm:ssZZZZZ" which is time.RFC3339
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

	slog.Debug("Sending HTTP request to GetUserProfileByID API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for GetUserProfileByID: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from GetUserProfileByID: %w", err)
	}
	slog.Debug("Received response from GetUserProfileByID API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	// Error handling based purely on StatusCode as per user feedback
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request GetUserProfileByID failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	// Attempt to unmarshal into the new wrapper structure
	var responseWrapper UserProfileAPIResponse
	if err := json.Unmarshal(responseBodyBytes, &responseWrapper); err != nil {
		// Log the body for debugging if unmarshal fails
		slog.Error("Failed to unmarshal GetUserProfileByID success response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, fmt.Errorf("failed to unmarshal GetUserProfileByID success response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	slog.Debug("Successfully fetched and parsed user profile.", "userID", userID, "nickname", responseWrapper.Data.Nickname)
	return &responseWrapper.Data, nil
}

// GetUserStats fetches a user's statistics by their UID.
func GetUserStats(userID string) (*UserStatsData, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty for GetUserStats")
	}
	apiURL := fmt.Sprintf("%s/v1/user-stats/get?uid=%s", constants.APIBaseURL, userID)
	slog.Debug("Fetching user stats by ID", "url", apiURL, "userID", userID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetUserStats: %w", err)
	}

	tm, err := GetTokenManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get token manager for GetUserStats: %w", err)
	}
	accessToken, err := tm.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token for GetUserStats: %w", err)
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

	slog.Debug("Sending HTTP request to GetUserStats API", "headers", req.Header)
	resp, err := GetHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for GetUserStats: %w", err)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from GetUserStats: %w", err)
	}
	slog.Debug("Received response from GetUserStats API", "statusCode", resp.StatusCode, "body", string(responseBodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request GetUserStats failed with status %d: %s", resp.StatusCode, string(responseBodyBytes))
	}

	var responseWrapper UserStatsAPIResponse // This type is defined in types.go
	if err := json.Unmarshal(responseBodyBytes, &responseWrapper); err != nil {
		slog.Error("Failed to unmarshal GetUserStats success response JSON", "error", err, "responseBody", string(responseBodyBytes))
		return nil, fmt.Errorf("failed to unmarshal GetUserStats success response JSON: %w. Body: %s", err, string(responseBodyBytes))
	}

	slog.Debug("Successfully fetched and parsed user stats.", "userID", userID)
	return &responseWrapper.Data, nil
}
