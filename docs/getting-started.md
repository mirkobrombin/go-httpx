# Getting Started

`go-httpx` wraps an `http.Client` transport stack and lets you compose middleware before optionally applying retry and circuit breaker policies from `go-foundation`.

## Building a client

- Use `httpx.New` to create a client with an optional `*http.Client` and any number of middleware values.
- Use `Header` to add standard headers to all outgoing requests.
- Use `WithRetry` and `WithBreaker` when you want to apply resiliency primitives around request execution.

## Middleware model

Middleware wraps `http.RoundTripper`, so it can be reused with existing transport-oriented patterns such as tracing, metrics, or authentication decorators.

## Error model

- Transport errors are returned directly.
- Retry delegates to `go-foundation/pkg/resiliency`.
- Circuit breaker failures surface the breaker error without masking the underlying behavior.
