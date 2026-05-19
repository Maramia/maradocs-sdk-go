package maradocs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxPollRetries = 120

// transport mirrors maradocs-sdk-ts FetchWrapper (JWT bearer, JSON, timeouts, pollResult).
type transport struct {
	jwt               string
	apiURLWithVersion string
	defaultTimeoutMs  *int
	httpClient        *http.Client
}

func newTransport(jwt, apiURLWithVersion string, defaultTimeoutMs *int) *transport {
	base := strings.TrimSuffix(apiURLWithVersion, "/")
	return &transport{
		jwt:               jwt,
		apiURLWithVersion: base,
		defaultTimeoutMs:  defaultTimeoutMs,
		httpClient:        http.DefaultClient,
	}
}

func (t *transport) url(path string) string {
	return t.apiURLWithVersion + path
}

func (t *transport) doHTTP(req *http.Request) (*http.Response, error) {
	return t.httpClient.Do(req)
}

func effectiveTimeout(requestMs *int, defaultMs *int) *int {
	if requestMs != nil {
		return requestMs
	}
	return defaultMs
}

func (t *transport) doJSON(
	ctx context.Context,
	method, path string,
	body any,
	out any,
	requestTimeoutMs *int,
) error {
	timeoutMs := effectiveTimeout(requestTimeoutMs, derefIntPtr(t.defaultTimeoutMs))
	var cancel context.CancelFunc
	if timeoutMs != nil {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*timeoutMs)*time.Millisecond)
		defer cancel()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, t.url(path), bodyReader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+t.jwt)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded && timeoutMs != nil {
			return fmt.Errorf("request timeout after %dms", *timeoutMs)
		}
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		if err := decodeSuccessBody(respBytes, out); err != nil {
			return fmt.Errorf("decode success body: %w (body=%s)", err, truncate(respBytes, 512))
		}
		return nil
	}

	var httpErr HttpErrorResponse
	if err := json.Unmarshal(respBytes, &httpErr); err == nil && httpErr.ApiErr.Name != "" {
		return &APIException{StatusCode: httpErr.StatusCode, Details: httpErr.ApiErr}
	}
	return fmt.Errorf("HTTP error: %s %s", resp.Status, truncate(respBytes, 256))
}

func (t *transport) getJSON(ctx context.Context, path string, out any, requestTimeoutMs *int) error {
	return t.doJSON(ctx, http.MethodGet, path, nil, out, requestTimeoutMs)
}

func (t *transport) postJSON(ctx context.Context, path string, body any, out any, requestTimeoutMs *int) error {
	return t.doJSON(ctx, http.MethodPost, path, body, out, requestTimeoutMs)
}

func (t *transport) putJSON(ctx context.Context, path string, body any, out any, requestTimeoutMs *int) error {
	return t.doJSON(ctx, http.MethodPut, path, body, out, requestTimeoutMs)
}

func (t *transport) patchJSON(ctx context.Context, path string, body any, out any, requestTimeoutMs *int) error {
	return t.doJSON(ctx, http.MethodPatch, path, body, out, requestTimeoutMs)
}

func (t *transport) deleteJSON(ctx context.Context, path string, body any, out any, requestTimeoutMs *int) error {
	return t.doJSON(ctx, http.MethodDelete, path, body, out, requestTimeoutMs)
}

func (t *transport) pollResult(ctx context.Context, path string, out any, requestTimeoutMs *int) error {
	timeoutMs := effectiveTimeout(requestTimeoutMs, derefIntPtr(t.defaultTimeoutMs))
	retries := 0
	for {
		var reqCtx context.Context
		var cancel context.CancelFunc
		if timeoutMs != nil {
			reqCtx, cancel = context.WithTimeout(ctx, time.Duration(*timeoutMs)*time.Millisecond)
		} else {
			reqCtx, cancel = context.WithCancel(ctx)
		}

		req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, t.url(path), nil)
		if err != nil {
			cancel()
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+t.jwt)

		resp, err := t.httpClient.Do(req)
		if err != nil {
			cancel()
			if reqCtx.Err() == context.DeadlineExceeded && timeoutMs != nil {
				return fmt.Errorf("request timeout after %dms", *timeoutMs)
			}
			return err
		}

		respBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel()
		if readErr != nil {
			return readErr
		}

		if resp.StatusCode == http.StatusOK {
			if err := decodeSuccessBody(respBytes, out); err != nil {
				return fmt.Errorf("poll decode: %w (body=%s)", err, truncate(respBytes, 512))
			}
			return nil
		}
		if resp.StatusCode == http.StatusAccepted {
			retries++
			if retries > maxPollRetries {
				return fmt.Errorf("failed processing the request in a reasonable time")
			}
			continue
		}

		var httpErr HttpErrorResponse
		if err := json.Unmarshal(respBytes, &httpErr); err == nil && httpErr.ApiErr.Name != "" {
			return &APIException{StatusCode: httpErr.StatusCode, Details: httpErr.ApiErr}
		}
		return fmt.Errorf("HTTP error during poll: %s %s", resp.Status, truncate(respBytes, 256))
	}
}

func derefIntPtr(p *int) *int {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

func decodeSuccessBody(body []byte, out any) error {
	if out == nil {
		return nil
	}
	s := bytes.TrimSpace(body)
	if len(s) == 0 || string(s) == "null" {
		return nil
	}
	return json.Unmarshal(body, out)
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
