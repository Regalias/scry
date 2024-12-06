package throttle

import (
	"net/http"
)

type RateLimitedHTTPTransport struct {
	limiter *HostRateLimiter
	// Underlying HTTP transport
	transport http.RoundTripper
}

// Returns a default HTTP RoundTripper with a flat TPS limit applied to http.DefaultTransport
func NewTransport(tps float64) http.RoundTripper {
	return NewWrappedTransport(http.DefaultTransport, tps)
}

// Returns a wrapped HTTP RoundTripper with a flat TPS limit applied
func NewWrappedTransport(transport http.RoundTripper, tps float64) http.RoundTripper {
	return &RateLimitedHTTPTransport{
		limiter:   NewHostRateLimiter(2),
		transport: transport,
	}
}

func (c *RateLimitedHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := c.limiter.GetLimiter(req.URL.Host).Wait(req.Context()); err != nil {
		return nil, err
	}
	return c.transport.RoundTrip(req)
}
