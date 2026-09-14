package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/agent/request"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8080", "secret-token")

	if client.baseURL != "http://localhost:8080" {
		t.Errorf("expected base url %q, got: %q", "http://localhost:8080", client.baseURL)
	}

	if client.token != "secret-token" {
		t.Errorf("expected token %q, got: %q", "secret-token", client.token)
	}

	if client.client == nil {
		t.Error("expected an http client, got nil")
	}
}

func TestSendMetricSuccess(t *testing.T) {
	var (
		gotMethod        string
		gotPath          string
		gotContentType   string
		gotAuthorization string
		gotBody          []byte
	)

	// fake server which only records the incoming request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		gotAuthorization = r.Header.Get("Authorization")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		gotBody = body

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret-token")

	metric := request.CreateMetricRequest{
		CPUUsage:       50.5,
		MemoryUsage:    60.5,
		DiskUsage:      70.5,
		NetworkReceive: 1000,
		NetworkSent:    2000,
	}

	err := client.SendMetric(metric)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method %s, got: %s", http.MethodPost, gotMethod)
	}

	if gotPath != "/agents/metrics" {
		t.Errorf("expected path %q, got: %q", "/agents/metrics", gotPath)
	}

	if gotContentType != "application/json" {
		t.Errorf("expected content type %q, got: %q", "application/json", gotContentType)
	}

	if gotAuthorization != "Bearer secret-token" {
		t.Errorf("expected authorization %q, got: %q", "Bearer secret-token", gotAuthorization)
	}

	var sentPayload request.CreateMetricRequest
	if err := json.Unmarshal(gotBody, &sentPayload); err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}

	if sentPayload != metric {
		t.Errorf("expected sent metric %+v, got: %+v", metric, sentPayload)
	}
}

func TestSendMetricUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret-token")

	err := client.SendMetric(request.CreateMetricRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "unexpected status code: 500") {
		t.Errorf("expected error to mention status code 500, got: %v", err)
	}
}

func TestSendMetricServerDown(t *testing.T) {
	// start and immediately close the server, so the connection is refused
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	server.Close()

	client := NewClient(server.URL, "secret-token")

	err := client.SendMetric(request.CreateMetricRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

