package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientLogin(t *testing.T) {
	t.Run("successful login stores the returned token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
			}
			if r.URL.Path != "/users/login" {
				t.Errorf("path = %q, want %q", r.URL.Path, "/users/login")
			}

			var req loginRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode request body: %v", err)
				return
			}
			if req.Username != "alice" {
				t.Errorf("username = %q, want %q", req.Username, "alice")
			}
			if req.Password != "s3cret" {
				t.Errorf("password = %q, want %q", req.Password, "s3cret")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"abc123"}`))
		}))
		defer server.Close()

		client := NewClient(server.URL)
		if err := client.Login(context.Background(), "alice", "s3cret"); err != nil {
			t.Fatalf("Login() error = %v, want nil", err)
		}

		if got := client.token; got != "abc123" {
			t.Errorf("token = %q, want %q", got, "abc123")
		}
	})

	t.Run("returns an error on a non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewClient(server.URL)
		err := client.Login(context.Background(), "alice", "s3cret")

		const want = "login failed with status code: 401"
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("Login() error = %v, want containing %q", err, want)
		}
		if client.token != "" {
			t.Errorf("token = %q, want empty", client.token)
		}
	})

	t.Run("returns an error when the response body is not valid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		client := NewClient(server.URL)
		if err := client.Login(context.Background(), "alice", "s3cret"); err == nil {
			t.Error("Login() error = nil, want a decode error")
		}
	})

	t.Run("returns an error when the request cannot be delivered", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		baseURL := server.URL
		server.Close()

		client := NewClient(baseURL)
		if err := client.Login(context.Background(), "alice", "s3cret"); err == nil {
			t.Error("Login() error = nil, want a network error")
		}
	})

	t.Run("returns an error when the context is already canceled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		client := NewClient(server.URL)
		if err := client.Login(ctx, "alice", "s3cret"); err == nil {
			t.Error("Login() error = nil, want a context error")
		}
	})
}
