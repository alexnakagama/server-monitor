package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// compile time interface assertion
// it assures in compile time that *metricServiceMock impliments the interface metricService
var _ metricService = (*metricServiceMock)(nil)

// creating a mock of the metric service
// we dont want the tests depend on the database
type metricServiceMock struct {
	calls atomic.Int64
	err   error
}

func (m *metricServiceMock) DeleteOldMetrics(ctx context.Context) error {
	m.calls.Add(1)
	if m.err != nil {
		return m.err
	}

	return nil
}

func (m *metricServiceMock) count() int {
	return int(m.calls.Load())
}

func TestRunWithCanceledContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mock := &metricServiceMock{}
		scheduler := NewScheduler(mock)

		// creates a context which we can cancel manually
		ctx, cancel := context.WithCancel(context.Background())
		// cancels the context
		cancel()

		// executes the scheduler with an already canceled context
		scheduler.Run(ctx)

		// verifies how many times the service called the scheduler
		if got := mock.count(); got != 0 {
			t.Errorf("DeleteOldMetrics called %d times, want 0", got)
		}
	})
}

func TestRunDeletesOldMetricsImmediatelyAndOnTicker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mock := &metricServiceMock{}
		scheduler := NewScheduler(mock)

		ctx, cancel := context.WithCancel(context.Background())

		go scheduler.Run(ctx)
		synctest.Wait()

		if got := mock.count(); got != 1 {
			t.Fatalf("DeleteOldMetrics called %d times after start, want 1", got)
		}

		time.Sleep(24 * time.Hour)
		synctest.Wait()

		if got := mock.count(); got != 2 {
			t.Fatalf("DeleteOldMetrics called %d times after 24h, want 2", got)
		}

		cancel()
		synctest.Wait()
	})
}

func TestRunKeepsRunningWhenDeleteOldMetricsFails(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		mock := &metricServiceMock{err: errors.New("boom")}
		scheduler := NewScheduler(mock)

		ctx, cancel := context.WithCancel(context.Background())

		go scheduler.Run(ctx)
		synctest.Wait()

		if got := mock.count(); got != 1 {
			t.Fatalf("DeleteOldMetrics called %d times after start, want 1", got)
		}

		time.Sleep(24 * time.Hour)
		synctest.Wait()

		if got := mock.count(); got != 2 {
			t.Fatalf("DeleteOldMetrics called %d times after 24h despite errors, want 2", got)
		}

		cancel()
		synctest.Wait()
	})
}
