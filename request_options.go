package maradocs

// RequestOptions can override per-request timeouts for poll-based endpoints.
type RequestOptions struct {
	// Timeout is the request timeout in milliseconds for this call (including poll steps).
	Timeout *int
}

func optTimeoutMs(opts *RequestOptions) *int {
	if opts == nil || opts.Timeout == nil {
		return nil
	}
	return opts.Timeout
}
