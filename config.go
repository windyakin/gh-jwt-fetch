package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	AppID          string
	PrivateKey     []byte
	InstallationID string
	RepoOwner      string
	RepoName       string
	FilePath       string
	OutputPath     string
	APIBaseURL     string
	Ref            string
	Interval       time.Duration
}

func LoadConfig() (Config, error) {
	var missing []string

	appID := os.Getenv("GH_APP_ID")
	if appID == "" {
		missing = append(missing, "GH_APP_ID")
	}

	installationID := os.Getenv("GH_APP_INSTALLATION_ID")
	if installationID == "" {
		missing = append(missing, "GH_APP_INSTALLATION_ID")
	}

	repoOwner := os.Getenv("GH_REPO_OWNER")
	if repoOwner == "" {
		missing = append(missing, "GH_REPO_OWNER")
	}

	repoName := os.Getenv("GH_REPO_NAME")
	if repoName == "" {
		missing = append(missing, "GH_REPO_NAME")
	}

	filePath := os.Getenv("GH_FILE_PATH")
	if filePath == "" {
		missing = append(missing, "GH_FILE_PATH")
	}

	outputPath := os.Getenv("GH_OUTPUT_PATH")
	if outputPath == "" {
		missing = append(missing, "GH_OUTPUT_PATH")
	}

	// Private key: direct content or file path
	var privateKey []byte
	if raw := os.Getenv("GH_APP_PRIVATE_KEY"); raw != "" {
		privateKey = []byte(raw)
	} else if keyPath := os.Getenv("GH_APP_PRIVATE_KEY_PATH"); keyPath != "" {
		data, err := os.ReadFile(keyPath)
		if err != nil {
			return Config{}, fmt.Errorf("failed to read private key file %s: %w", keyPath, err)
		}
		privateKey = data
	} else {
		missing = append(missing, "GH_APP_PRIVATE_KEY or GH_APP_PRIVATE_KEY_PATH")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	apiBaseURL := os.Getenv("GH_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.github.com"
	}
	apiBaseURL = strings.TrimRight(apiBaseURL, "/")

	var interval time.Duration
	if intervalStr := os.Getenv("GH_INTERVAL"); intervalStr != "" {
		var err error
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			return Config{}, fmt.Errorf("invalid GH_INTERVAL %q: %w", intervalStr, err)
		}
		if interval < 1*time.Second {
			return Config{}, fmt.Errorf("GH_INTERVAL must be at least 1s, got %s", interval)
		}
	}

	return Config{
		AppID:          appID,
		PrivateKey:     privateKey,
		InstallationID: installationID,
		RepoOwner:      repoOwner,
		RepoName:       repoName,
		FilePath:       filePath,
		OutputPath:     outputPath,
		APIBaseURL:     apiBaseURL,
		Ref:            os.Getenv("GH_REF"),
		Interval:       interval,
	}, nil
}
