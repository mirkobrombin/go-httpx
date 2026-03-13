package httpx

import (
	"net/http"
	"time"

	"github.com/mirkobrombin/go-foundation/pkg/resiliency"
)

// Middleware wraps an http.RoundTripper.
type Middleware func(next http.RoundTripper) http.RoundTripper

// Client is an extensible HTTP client with middleware and optional resiliency hooks.
type Client struct {
	client *http.Client
	rt     http.RoundTripper
	mw     []Middleware

	retryOpts []func(*resiliency.RetryOptions)
	breaker   *resiliency.CircuitBreaker
}

// New creates a client with optional middleware.
func New(c *http.Client, mw ...Middleware) *Client {
	if c == nil {
		c = &http.Client{Timeout: 15 * time.Second}
	}

	rt := c.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}

	return &Client{client: c, rt: rt, mw: mw}
}

// WithRetry attaches retry options to the client.
func (c *Client) WithRetry(opts ...func(*resiliency.RetryOptions)) *Client {
	c.retryOpts = opts
	return c
}

// WithBreaker attaches a circuit breaker to the client.
func (c *Client) WithBreaker(b *resiliency.CircuitBreaker) *Client {
	c.breaker = b
	return c
}

func (c *Client) buildTransport() http.RoundTripper {
	rt := c.rt
	for i := len(c.mw) - 1; i >= 0; i-- {
		rt = c.mw[i](rt)
	}
	return rt
}

// Do sends a request applying middleware and optional resiliency behaviors.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	rt := c.buildTransport()
	var resp *http.Response

	fn := func() error {
		r, err := rt.RoundTrip(req)
		resp = r
		return err
	}

	if c.breaker != nil {
		if len(c.retryOpts) > 0 {
			return resp, c.breaker.Execute(func() error {
				return resiliency.Retry(req.Context(), fn, c.retryOpts...)
			})
		}
		return resp, c.breaker.Execute(fn)
	}

	if len(c.retryOpts) > 0 {
		return resp, resiliency.Retry(req.Context(), fn, c.retryOpts...)
	}

	return resp, fn()
}
