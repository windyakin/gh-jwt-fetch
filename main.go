package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	client := &GitHubClient{
		BaseURL:    cfg.APIBaseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}

	if err := fetch(cfg, client); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	if cfg.Interval == 0 {
		return
	}

	// Loop mode
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	log.Printf("running every %s (next at %s)", cfg.Interval, time.Now().Add(cfg.Interval).Format(time.RFC3339))

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-ticker.C:
			if err := fetch(cfg, client); err != nil {
				log.Printf("error: %v", err)
			}
		}
	}
}

func fetch(cfg Config, client *GitHubClient) error {
	jwt, err := GenerateJWT(cfg.AppID, cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("JWT generation error: %w", err)
	}

	token, err := client.GetInstallationToken(jwt, cfg.InstallationID)
	if err != nil {
		return fmt.Errorf("installation token error: %w", err)
	}

	content, err := client.DownloadFile(token, cfg.RepoOwner, cfg.RepoName, cfg.FilePath, cfg.Ref)
	if err != nil {
		return fmt.Errorf("file download error: %w", err)
	}

	if dir := filepath.Dir(cfg.OutputPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	if err := os.WriteFile(cfg.OutputPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Printf("downloaded %s/%s:%s to %s", cfg.RepoOwner, cfg.RepoName, cfg.FilePath, cfg.OutputPath)
	return nil
}
