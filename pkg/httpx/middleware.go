package httpx

import "net/http"

// Header adds a header to all outgoing requests.
func Header(key, value string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			cloned := cloneRequest(r)
			cloned.Header.Set(key, value)
			return next.RoundTrip(cloned)
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func cloneRequest(r *http.Request) *http.Request {
	cloned := r.Clone(r.Context())
	cloned.Header = make(http.Header, len(r.Header))
	for k, v := range r.Header {
		cloned.Header[k] = append([]string(nil), v...)
	}
	return cloned
}
