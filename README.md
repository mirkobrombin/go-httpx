# Go HTTPX

> [!CAUTION]
> go-httpx is now part of the [go-foundation](https://github.com/mirkobrombin/go-foundation) framework. The v1.0.0 release mirrors go-httpx v0.2.0, but future versions may introduce breaking changes. Please migrate your project.

A small **middleware-friendly** HTTP client for Go with optional **retry** and **circuit breaker** integration.

## Features

- **Composable Middleware:** Build request pipelines on top of `http.RoundTripper`.
- **Retry Integration:** Reuse `go-foundation` retry primitives when requests fail.
- **Circuit Breaker Support:** Protect upstream dependencies with a shared breaker.
- **Small API Surface:** Keep client construction and request execution easy to test.

## Installation

```bash
go get github.com/mirkobrombin/go-httpx
```

## Quick Start

```go
package main

import (
    "net/http"

    "github.com/mirkobrombin/go-foundation/pkg/resiliency"
    "github.com/mirkobrombin/go-httpx/pkg/httpx"
)

func main() {
    client := httpx.New(nil, httpx.Header("User-Agent", "my-client/1.0"))
    client.WithRetry(resiliency.WithAttempts(5))

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
```

## Documentation

- [Getting Started](docs/getting-started.md)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
