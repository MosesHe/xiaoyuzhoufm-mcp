package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"xiaoyuzhoufm-mcp/internal/server" // Import the server package
	"xiaoyuzhoufm-mcp/internal/xyzclient"

	"github.com/lmittmann/tint"
)

const (
	defaultAreaCode         = "+86"
	maxVerificationAttempts = 3
)

func main() {
	// Initialize structured logger
	tintOptions := &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.DateTime,
		AddSource:  true,
	}
	handler := tint.NewHandler(os.Stderr, tintOptions) // Log to Stderr

	logger := slog.New(handler)
	slog.SetDefault(logger)

	if len(os.Args) > 1 && os.Args[1] == "init" {
		slog.Debug("Running in init mode for interactive login.")
		tm, err := xyzclient.GetTokenManager()
		if err != nil {
			slog.Error("Failed to get token manager for init.", "error", err)
			fmt.Printf("Error initializing for login: %v\n", err)
			os.Exit(1)
		}
		interactiveLogin(tm) // Call the new combined interactiveLogin function
		slog.Debug("Initialization complete. Token saved. Exiting.")
		os.Exit(0)
	} else {
		// Default server mode
		slog.Debug("MCP Server starting in default mode...")
		tm, err := xyzclient.GetTokenManager()
		if err != nil {
			slog.Error("Failed to get token manager.", "error", err)
			os.Exit(1)
		}

		userTokenPath, pathErr := xyzclient.GetUserTokenPath()
		if pathErr != nil {
			slog.Error("Failed to determine user token path.", "error", pathErr)
			os.Exit(1)
		}

		slog.Debug("Attempting to load token from user path.", "path", userTokenPath)
		if err := tm.LoadTokenFromPath(userTokenPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				slog.Error("Token file not found at user path. Please run './xiaoyuzhoufm-mcp init' first.", "path", userTokenPath)
				fmt.Fprintln(os.Stderr, "Error: Token not found. Please run './xiaoyuzhoufm-mcp init' to login and create the token.")
			} else {
				slog.Error("Failed to load token from user path. The token file might be corrupted.", "path", userTokenPath, "error", err)
				fmt.Fprintf(os.Stderr, "Error: Failed to load token from %s. It might be corrupted. Try running './xiaoyuzhoufm-mcp init' again.\n", userTokenPath)
			}
			os.Exit(1)
		}
		slog.Debug("Token loaded successfully from user path.")
		server.RunStdioServer()
	}
	slog.Debug("MCP Server closed.")
}

// interactiveLogin handles the full interactive login process,
// populates the TokenManager, and saves the token to the user-specific path.
func interactiveLogin(tm *xyzclient.TokenManager) {
	slog.Debug("Starting interactive login and token save process.")

	reader := bufio.NewReader(os.Stdin)
	var areaCode string
	for {
		fmt.Printf("Enter your area code (e.g., +86, press Enter for default %s): ", defaultAreaCode)
		input, _ := reader.ReadString('\n')
		areaCode = strings.TrimSpace(input)
		if areaCode == "" {
			areaCode = defaultAreaCode
		}
		if isValidAreaCode(areaCode) {
			break
		}
		slog.Warn("Invalid area code format. Please try again.", "input", areaCode)
		fmt.Println("Invalid area code format. It should start with '+' followed by 1 to 3 digits (e.g., +86).")
	}

	var phoneNumber string
	for {
		fmt.Print("Enter your phone number (digits only): ")
		input, _ := reader.ReadString('\n')
		phoneNumber = strings.TrimSpace(input)
		if isValidPhoneNumber(phoneNumber) {
			break
		}
		slog.Warn("Invalid phone number format. Please enter digits only.", "input", phoneNumber)
		fmt.Println("Invalid phone number format. Please enter 7 to 15 digits.")
	}

	slog.Debug("Requesting verification code.", "areaCode", areaCode, "phoneNumber", phoneNumber)
	if err := xyzclient.RequestVerificationCode(areaCode, phoneNumber); err != nil {
		slog.Error("Error requesting verification code.", "error", err)
		fmt.Printf("Error requesting verification code: %v\n", err)
		os.Exit(1)
	}
	slog.Debug("Verification code request sent. Please check your phone.")
	fmt.Println("Verification code request sent. Please check your phone.")

	var verificationCode string
	loginAttempts := 0
	loginSuccessful := false
	for loginAttempts < maxVerificationAttempts {
		fmt.Printf("Enter the 4-digit verification code (attempt %d/%d): ", loginAttempts+1, maxVerificationAttempts)
		input, _ := reader.ReadString('\n')
		verificationCode = strings.TrimSpace(input)

		if !isValidVerificationCode(verificationCode) {
			slog.Warn("Invalid verification code format. It must be 4 digits.", "input", verificationCode)
			fmt.Println("Invalid verification code format. It must be 4 digits. Please try again.")
			continue
		}

		slog.Debug("Attempting to login with verification code.", "attempt", loginAttempts+1)
		accessToken, refreshToken, uid, nickname, err := xyzclient.LoginWithCode(areaCode, phoneNumber, verificationCode)
		if err == nil {
			slog.Debug("Login successful.", "uid", uid, "nickname", nickname)
			tm.AccessToken = accessToken
			tm.RefreshToken = refreshToken
			tm.Uid = uid
			tm.Nickname = nickname
			loginSuccessful = true
			break
		}
		slog.Warn("LoginWithCode failed.", "error", err, "attempt", loginAttempts+1)
		fmt.Printf("Login failed: %v. ", err)
		loginAttempts++
		if loginAttempts < maxVerificationAttempts {
			fmt.Printf("Please try again. %d attempts remaining.\n", maxVerificationAttempts-loginAttempts)
		} else {
			fmt.Println("Maximum login attempts reached. Exiting.")
			os.Exit(1)
		}
	}

	if !loginSuccessful {
		slog.Error("Interactive login steps did not complete successfully.")
		fmt.Println("Login process did not complete successfully. Exiting.")
		os.Exit(1)
	}

	// Save logic
	userTokenPath, pathErr := xyzclient.GetUserTokenPath()
	if pathErr != nil {
		slog.Error("Failed to get user token path for saving.", "error", pathErr)
		fmt.Printf("Error determining where to save token: %v\n", pathErr)
		os.Exit(1)
	}

	slog.Debug("Attempting to save token to user path.", "path", userTokenPath)
	if err := tm.SaveTokenToPath(userTokenPath); err != nil {
		slog.Error("Failed to save token to user path after successful login.", "path", userTokenPath, "error", err)
		fmt.Printf("Error: Successfully logged in, but failed to save token to %s: %v\n", userTokenPath, err)
		fmt.Println("Please check permissions or disk space and try './xiaoyuzhoufm-mcp init' again.")
		os.Exit(1)
	}

	slog.Debug("Token saved successfully to user path.", "path", userTokenPath)
	fmt.Println("Login successful and token saved!")
}

func isValidAreaCode(areaCode string) bool {
	re := regexp.MustCompile(`^\+\d{1,3}$`)
	return re.MatchString(areaCode)
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\d{7,15}$`)
	return re.MatchString(phone)
}

func isValidVerificationCode(code string) bool {
	re := regexp.MustCompile(`^\d{4}$`)
	return re.MatchString(code)
}
