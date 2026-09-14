package agent

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// compile time interface assertion
// it assures in compile time that *metricRoundTripper implements http.RoundTripper
var _ http.RoundTripper = (*metricRoundTripper)(nil)

// metricRoundTripper is an http.RoundTripper mock assigned to a real *Client,
// so Agent.Run can be tested without any network call.
// It records every request and answers with a configurable status code.
type metricRoundTripper struct {
	calls  atomic.Int64
	status int
}

func (m *metricRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	m.calls.Add(1)

	return &http.Response{
		StatusCode: m.status,
		Status:     http.StatusText(m.status),
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func (m *metricRoundTripper) count() int {
	return int(m.calls.Load())
}

// newSpyClient returns a real Client backed by a mock transport,
// so tests can observe every request and force SendMetric to fail.
func newSpyClient(status int) (*Client, *metricRoundTripper) {
	client := NewClient("http://localhost", "secret-token")
	transport := &metricRoundTripper{status: status}
	client.client = &http.Client{Transport: transport}

	return client, transport
}

// waitForSendCount blocks until SendMetric has been called at least want
// times, failing the test if that never happens before the timeout.
func waitForSendCount(t *testing.T, transport *metricRoundTripper, want int, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for transport.count() < want {
		if time.Now().After(deadline) {
			t.Fatalf("SendMetric called %d times, want %d (timed out)", transport.count(), want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRunWithCanceledContext(t *testing.T) {
	client, transport := newSpyClient(http.StatusCreated)
	agent := New(time.Hour, client)

	// creates a context which we can cancel manually
	ctx, cancel := context.WithCancel(context.Background())
	// cancels the context
	cancel()

	go agent.Run(ctx)

	// the agent already started collecting the first metric, so it still
	// sends that one before noticing the canceled context and stopping
	waitForSendCount(t, transport, 1, 10*time.Second)
}

func TestRunStopsWhenContextIsCanceled(t *testing.T) {
	client, transport := newSpyClient(http.StatusCreated)
	agent := New(time.Hour, client)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	go func() {
		agent.Run(ctx)
		close(done)
	}()

	waitForSendCount(t, transport, 1, 10*time.Second)

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("agent did not stop after context cancellation")
	}
}

func TestRunSendsMetricsOnEachTick(t *testing.T) {
	// the metric collection takes about 1.5 seconds, so a 3 second interval
	// is large enough to let each tick drive a whole collection round
	client, transport := newSpyClient(http.StatusCreated)
	agent := New(3*time.Second, client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.Run(ctx)

	// the first metric is collected and sent right after the agent starts
	waitForSendCount(t, transport, 1, 10*time.Second)

	// after the first tick interval a second metric must be sent
	waitForSendCount(t, transport, 2, 10*time.Second)
}

func TestRunStopsWhenSendingFails(t *testing.T) {
	// the server always answers with an error status code
	client, transport := newSpyClient(http.StatusInternalServerError)
	agent := New(time.Millisecond, client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.Run(ctx)

	// the first send attempt fails and makes Run return
	waitForSendCount(t, transport, 1, 10*time.Second)

	// wait longer than a collection round: if the loop had kept running,
	// we would see another request
	time.Sleep(2 * time.Second)

	if got := transport.count(); got != 1 {
		t.Errorf("SendMetric called %d times after failure, want 1", got)
	}
}

