package maradocs

import (
	"context"
	"net/http"
	"time"
)

// HealthcheckEp implements unauthenticated health probes (GET /healthcheck/ping).
type HealthcheckEp struct {
	baseURL string
	client  *http.Client
	timeout *int
}

func newHealthcheckEp(apiURLWithVersion string, timeoutMs *int) *HealthcheckEp {
	t := newTransport("", apiURLWithVersion, timeoutMs)
	return &HealthcheckEp{
		baseURL: t.apiURLWithVersion,
		client:  t.httpClient,
		timeout: timeoutMs,
	}
}

// Ping returns true if the health endpoint responds with HTTP 2xx.
func (h *HealthcheckEp) Ping(ctx context.Context) (bool, error) {
	if h.timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*h.timeout)*time.Millisecond)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseURL+"/healthcheck/ping", nil)
	if err != nil {
		return false, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300, nil
}
