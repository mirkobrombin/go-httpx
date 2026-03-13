package httpx_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mirkobrombin/go-foundation/pkg/resiliency"
	"github.com/mirkobrombin/go-httpx/pkg/httpx"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestHeaderMiddlewareAddsHeaderWithoutMutatingRequest(t *testing.T) {
	originalHeader := http.Header{}
	originalHeader.Set("X-Original", "1")

	var seenHeader string
	client := httpx.New(&http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		seenHeader = r.Header.Get("X-Test")
		if r.Header.Get("X-Original") != "1" {
			t.Fatalf("RoundTrip() original header missing")
		}
		return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
	})}, httpx.Header("X-Test", "1"))

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}
	req.Header = originalHeader

	if _, err := client.Do(req); err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if seenHeader != "1" {
		t.Fatalf("RoundTrip() saw header %q, want %q", seenHeader, "1")
	}
	if req.Header.Get("X-Test") != "" {
		t.Fatalf("request header mutated unexpectedly")
	}
}

func TestClientWithRetryRetriesTransportErrors(t *testing.T) {
	attempts := 0
	client := httpx.New(&http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("temporary")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok"))}, nil
	})})
	client.WithRetry(
		resiliency.WithAttempts(3),
		resiliency.WithDelay(time.Nanosecond, time.Nanosecond),
	)

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if attempts != 3 {
		t.Fatalf("Do() attempts = %d, want %d", attempts, 3)
	}
}

func TestClientWithBreakerStopsOpenCircuit(t *testing.T) {
	breaker := resiliency.NewCircuitBreaker(1, time.Hour)
	client := httpx.New(&http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})})
	client.WithBreaker(breaker)

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	if _, err := client.Do(req); err == nil {
		t.Fatalf("first Do() error = nil, want transport error")
	}

	if _, err := client.Do(req); !errors.Is(err, resiliency.ErrCircuitOpen) {
		t.Fatalf("second Do() error = %v, want ErrCircuitOpen", err)
	}
}
