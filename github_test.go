package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInstallationToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/app/installations/12345/access_tokens" {
			t.Errorf("path = %s, want /app/installations/12345/access_tokens", r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-jwt" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-jwt")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"token":"inst-token-abc","expires_at":"2099-01-01T00:00:00Z"}`))
	}))
	defer server.Close()

	client := &GitHubClient{BaseURL: server.URL, HTTPClient: server.Client()}
	token, err := client.GetInstallationToken("test-jwt", "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "inst-token-abc" {
		t.Errorf("token = %q, want %q", token, "inst-token-abc")
	}
}

func TestGetInstallationToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Bad credentials"}`))
	}))
	defer server.Close()

	client := &GitHubClient{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := client.GetInstallationToken("bad-jwt", "12345")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDownloadFile_Success(t *testing.T) {
	fileContent := []byte("hello world content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/repos/owner/repo/contents/path/to/file.txt" {
			t.Errorf("path = %s, want /repos/owner/repo/contents/path/to/file.txt", r.URL.Path)
		}
		if ref := r.URL.Query().Get("ref"); ref != "main" {
			t.Errorf("ref = %q, want %q", ref, "main")
		}
		accept := r.Header.Get("Accept")
		if accept != "application/vnd.github.raw+json" {
			t.Errorf("Accept = %q, want %q", accept, "application/vnd.github.raw+json")
		}
		w.WriteHeader(http.StatusOK)
		w.Write(fileContent)
	}))
	defer server.Close()

	client := &GitHubClient{BaseURL: server.URL, HTTPClient: server.Client()}
	data, err := client.DownloadFile("token", "owner", "repo", "path/to/file.txt", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(fileContent) {
		t.Errorf("content = %q, want %q", string(data), string(fileContent))
	}
}

func TestDownloadFile_NoRef(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ref := r.URL.Query().Get("ref"); ref != "" {
			t.Errorf("ref should be empty, got %q", ref)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("content"))
	}))
	defer server.Close()

	client := &GitHubClient{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := client.DownloadFile("token", "owner", "repo", "file.txt", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDownloadFile_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	client := &GitHubClient{BaseURL: server.URL, HTTPClient: server.Client()}
	_, err := client.DownloadFile("token", "owner", "repo", "nonexistent", "")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
