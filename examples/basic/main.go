package main

import (
	"net/http"
	"time"

	"github.com/mirkobrombin/go-foundation/pkg/resiliency"
	"github.com/mirkobrombin/go-httpx/pkg/httpx"
)

func main() {
	client := httpx.New(nil, httpx.Header("User-Agent", "go-httpx-example/1.0"))
	client.WithRetry(
		resiliency.WithAttempts(3),
		resiliency.WithDelay(50*time.Millisecond, 100*time.Millisecond),
	)

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
