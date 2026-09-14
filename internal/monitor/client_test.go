package monitor

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGetServers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			if r.URL.Path != "/servers" {
				t.Errorf("expected path /servers, got %s", r.URL.Path)
			}

			auth := r.Header.Get("Authorization")
			if auth != "Bearer abc123" {
				t.Errorf("expected Authorization Bearer abc123, got %s", auth)
			}

			w.Header().Set("Content-Type", "application/json")

			fmt.Fprint(w, `[
				{
					"id": 1,
					"name": "server-1",
					"hostname": "server-1.local",
					"os": "linux",
					"created_at": "2026-09-14T12:00:00Z"
				}
			]`)
		}))
		defer server.Close()

		client := NewClient(server.URL)
		client.token = "abc123"

		servers, err := client.GetServers(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(servers) != 1 {
			t.Fatalf("expected 1 server, got %d", len(servers))
		}

		if servers[0].ID != 1 {
			t.Errorf("expected ID 1, got %d", servers[0].ID)
		}

		if servers[0].Name != "server-1" {
			t.Errorf("expected name server-1, got %s", servers[0].Name)
		}
	})
}
