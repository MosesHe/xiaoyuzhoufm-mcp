package xyzclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	tokenFileName = "token.json"
)

type TokenManager struct {
	AccessToken          string `json:"access_token"`
	RefreshToken         string `json:"refresh_token"`
	Uid                  string `json:"uid"`
	Nickname             string `json:"nickname"`
	LastUpdatedTimestamp int64  `json:"last_updated_timestamp,omitempty"`
	loadedTokenPath      string `json:"-"` // Path from which token was loaded or to which it was last saved. Not persisted in JSON.
}

var (
	instance         *TokenManager
	tokenManagerOnce sync.Once
	initErr          error
)

// GetUserTokenPath returns the OS-specific path for storing the token in the user's home directory.
// Path is typically ~/.mcp/xiaoyuzhoufm-mcp/token.json
func GetUserTokenPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	tokenDir := filepath.Join(homeDir, ".mcp", "xiaoyuzhoufm-mcp")
	return filepath.Join(tokenDir, tokenFileName), nil
}

func GetTokenManager() (*TokenManager, error) {
	tokenManagerOnce.Do(func() {
		slog.Debug("TokenManager instance being created.")
		instance = &TokenManager{}
	})
	return instance, initErr // initErr might be nil or set by a previous critical error if we decide to.
}

// LoadTokenFromPath loads token data from the specified file path.
func (tm *TokenManager) LoadTokenFromPath(tokenPath string) error {
	if tokenPath == "" {
		return fmt.Errorf("token path cannot be empty")
	}
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fs.ErrNotExist // Return fs.ErrNotExist specifically so callers can check for it.
		}
		return fmt.Errorf("failed to read token file %s: %w", tokenPath, err)
	}

	if err := json.Unmarshal(data, tm); err != nil {
		return fmt.Errorf("failed to unmarshal token data from %s: %w", tokenPath, err)
	}

	if tm.AccessToken == "" || tm.RefreshToken == "" {
		slog.Warn("Loaded token data from path is incomplete", "path", tokenPath)
		return fmt.Errorf("incomplete token data in %s", tokenPath)
	}
	tm.loadedTokenPath = tokenPath // Store the path from which the token was loaded
	slog.Debug("Token loaded successfully from path.", "path", tokenPath)
	return nil
}

// SaveTokenToPath saves the current token data to the specified file path.
// It creates necessary directories if they don't exist.
func (tm *TokenManager) SaveTokenToPath(tokenPath string) error {
	if tokenPath == "" {
		return fmt.Errorf("token path cannot be empty for saving")
	}
	now := time.Now().Unix()
	tm.LastUpdatedTimestamp = now

	data, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(tokenPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write token file %s: %w", tokenPath, err)
	}
	tm.loadedTokenPath = tokenPath // Update the path to which token was last saved
	slog.Debug("Token saved successfully to path.", "path", tokenPath)
	return nil
}

func (tm *TokenManager) GetAccessToken() (string, error) {
	if tm.AccessToken == "" {
		slog.Warn("GetAccessToken called but access token is empty (initial state or previous error).")
		return "", fmt.Errorf("not authenticated: access token is empty")
	}

	const tokenTimeoutSeconds = 20 * 60 // 20 minutes
	currentTime := time.Now().Unix()

	if tm.LastUpdatedTimestamp > 0 && (currentTime-tm.LastUpdatedTimestamp > tokenTimeoutSeconds) {
		slog.Debug("Access token may have expired, attempting to refresh.", "lastUpdated", time.Unix(tm.LastUpdatedTimestamp, 0), "currentTime", time.Unix(currentTime, 0))
		err := tm.RefreshAccessToken() // RefreshAccessToken will save the token if successful
		if err != nil {
			return "", fmt.Errorf("failed to refresh token, authentication may be required: %w", err)
		}
		slog.Debug("Access token refreshed successfully during GetAccessToken.")
		// After successful refresh, ensure accessToken is not empty before returning
		if tm.AccessToken == "" {
			return "", fmt.Errorf("authentication failed: token became empty after refresh")
		}
	}

	return tm.AccessToken, nil
}

// RefreshAccessToken attempts to refresh the access token using the refresh token.
// If successful, it updates the AccessToken and RefreshToken fields and saves the token.
func (tm *TokenManager) RefreshAccessToken() error {
	slog.Debug("Attempting to refresh access token")
	if tm.RefreshToken == "" {
		return fmt.Errorf("cannot refresh token: refresh token is empty")
	}

	newAccessToken, newRefreshToken, err := PerformTokenRefresh(tm.RefreshToken)
	if err != nil {
		return fmt.Errorf("PerformTokenRefresh failed: %w", err)
	}

	tm.AccessToken = newAccessToken
	tm.RefreshToken = newRefreshToken
	// Update timestamp upon successful refresh before saving
	tm.LastUpdatedTimestamp = time.Now().Unix()
	slog.Debug("Access token and refresh token updated locally after refresh.")

	if tm.loadedTokenPath == "" {
		slog.Warn("Cannot persist refreshed token: loadedTokenPath is not set. Token refreshed in memory only.")
		return nil // Token is refreshed in memory.
	}

	if err := tm.SaveTokenToPath(tm.loadedTokenPath); err != nil {
		return fmt.Errorf("successfully refreshed token, but failed to save to file %s: %w", tm.loadedTokenPath, err)
	}
	slog.Debug("Refreshed token saved successfully.", "path", tm.loadedTokenPath)
	return nil
}
