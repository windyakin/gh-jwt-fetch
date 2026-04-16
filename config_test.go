package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_AllRequired(t *testing.T) {
	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "path/to/file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AppID != "12345" {
		t.Errorf("AppID = %q, want %q", cfg.AppID, "12345")
	}
	if cfg.InstallationID != "67890" {
		t.Errorf("InstallationID = %q, want %q", cfg.InstallationID, "67890")
	}
	if cfg.RepoOwner != "owner" {
		t.Errorf("RepoOwner = %q, want %q", cfg.RepoOwner, "owner")
	}
	if cfg.RepoName != "repo" {
		t.Errorf("RepoName = %q, want %q", cfg.RepoName, "repo")
	}
	if cfg.FilePath != "path/to/file" {
		t.Errorf("FilePath = %q, want %q", cfg.FilePath, "path/to/file")
	}
	if cfg.OutputPath != "/tmp/out" {
		t.Errorf("OutputPath = %q, want %q", cfg.OutputPath, "/tmp/out")
	}
	if cfg.APIBaseURL != "https://api.github.com" {
		t.Errorf("APIBaseURL = %q, want %q", cfg.APIBaseURL, "https://api.github.com")
	}
	if cfg.Interval != 0 {
		t.Errorf("Interval = %v, want 0 (oneshot)", cfg.Interval)
	}
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	// Clear all env vars
	for _, key := range []string{"GH_APP_ID", "GH_APP_PRIVATE_KEY", "GH_APP_PRIVATE_KEY_PATH", "GH_APP_INSTALLATION_ID", "GH_REPO_OWNER", "GH_REPO_NAME", "GH_FILE_PATH", "GH_OUTPUT_PATH"} {
		t.Setenv(key, "")
	}

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadConfig_PrivateKeyFromFile(t *testing.T) {
	keyContent := "-----BEGIN RSA PRIVATE KEY-----\nfromfile\n-----END RSA PRIVATE KEY-----"
	keyFile := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(keyFile, []byte(keyContent), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "")
	t.Setenv("GH_APP_PRIVATE_KEY_PATH", keyFile)
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(cfg.PrivateKey) != keyContent {
		t.Errorf("PrivateKey = %q, want %q", string(cfg.PrivateKey), keyContent)
	}
}

func TestLoadConfig_PrivateKeyDirectTakesPrecedence(t *testing.T) {
	directKey := "-----BEGIN RSA PRIVATE KEY-----\ndirect\n-----END RSA PRIVATE KEY-----"

	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", directKey)
	t.Setenv("GH_APP_PRIVATE_KEY_PATH", "/nonexistent/path")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(cfg.PrivateKey) != directKey {
		t.Errorf("PrivateKey = %q, want %q", string(cfg.PrivateKey), directKey)
	}
}

func TestLoadConfig_APIBaseURLTrailingSlash(t *testing.T) {
	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")
	t.Setenv("GH_API_BASE_URL", "https://ghe.example.com/api/v3/")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIBaseURL != "https://ghe.example.com/api/v3" {
		t.Errorf("APIBaseURL = %q, want %q", cfg.APIBaseURL, "https://ghe.example.com/api/v3")
	}
}

func TestLoadConfig_IntervalValid(t *testing.T) {
	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")
	t.Setenv("GH_INTERVAL", "5m")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 5*time.Minute {
		t.Errorf("Interval = %v, want 5m", cfg.Interval)
	}
}

func TestLoadConfig_IntervalInvalid(t *testing.T) {
	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")
	t.Setenv("GH_INTERVAL", "invalid")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestLoadConfig_IntervalTooShort(t *testing.T) {
	t.Setenv("GH_APP_ID", "12345")
	t.Setenv("GH_APP_PRIVATE_KEY", "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----")
	t.Setenv("GH_APP_INSTALLATION_ID", "67890")
	t.Setenv("GH_REPO_OWNER", "owner")
	t.Setenv("GH_REPO_NAME", "repo")
	t.Setenv("GH_FILE_PATH", "file")
	t.Setenv("GH_OUTPUT_PATH", "/tmp/out")
	t.Setenv("GH_INTERVAL", "500ms")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for interval < 1s")
	}
}
