package maradocs

import (
	"context"
	"fmt"
	"net/url"
)

// VideoEp implements /video/* endpoints.
type VideoEp struct {
	t *transport
}

func newVideoEp(t *transport) *VideoEp {
	return &VideoEp{t: t}
}

// Validate validates and transcodes video.
func (e *VideoEp) Validate(ctx context.Context, req VideoValidateRequest, opts *RequestOptions) (*VideoValidateResponse, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, "/video/validate", req, &task, to); err != nil {
		return nil, err
	}
	var out VideoValidateResponse
	if err := e.t.pollResult(ctx, fmt.Sprintf("/video/validate/%s", url.PathEscape(task.JobID)), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}

// AudioEp implements /audio/* endpoints.
type AudioEp struct {
	t *transport
}

func newAudioEp(t *transport) *AudioEp {
	return &AudioEp{t: t}
}

// Validate validates and transcodes audio.
func (e *AudioEp) Validate(ctx context.Context, req AudioValidateRequest, opts *RequestOptions) (*AudioValidateResponse, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, "/audio/validate", req, &task, to); err != nil {
		return nil, err
	}
	var out AudioValidateResponse
	if err := e.t.pollResult(ctx, fmt.Sprintf("/audio/validate/%s", url.PathEscape(task.JobID)), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}
