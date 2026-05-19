package maradocs

import (
	"context"
	"fmt"
	"net/url"
)

// HtmlEp implements /html/* endpoints.
type HtmlEp struct {
	t *transport
}

func newHtmlEp(t *transport) *HtmlEp {
	return &HtmlEp{t: t}
}

func (e *HtmlEp) postPoll(ctx context.Context, postPath string, req any, pollFmt string, out any, opts *RequestOptions) error {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, postPath, req, &task, to); err != nil {
		return err
	}
	return e.t.pollResult(ctx, fmt.Sprintf(pollFmt, url.PathEscape(task.JobID)), out, to)
}

// Validate validates HTML.
func (e *HtmlEp) Validate(ctx context.Context, req HtmlValidateRequest, opts *RequestOptions) (*HtmlValidateResponse, error) {
	var out HtmlValidateResponse
	if err := e.postPoll(ctx, "/html/validate", req, "/html/validate/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToPdf converts HTML to PDF.
func (e *HtmlEp) ToPdf(ctx context.Context, req HtmlToPdfRequest, opts *RequestOptions) (*HtmlToPdfResponse, error) {
	var out HtmlToPdfResponse
	if err := e.postPoll(ctx, "/html/to/pdf", req, "/html/to/pdf/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// EmailEp implements /email/* endpoints.
type EmailEp struct {
	t *transport
}

func newEmailEp(t *transport) *EmailEp {
	return &EmailEp{t: t}
}

// Validate validates an email file.
func (e *EmailEp) Validate(ctx context.Context, req EmailValidateRequest, opts *RequestOptions) (*EmailValidateResponse, error) {
	var out EmailValidateResponse
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, "/email/validate", req, &task, to); err != nil {
		return nil, err
	}
	if err := e.t.pollResult(ctx, fmt.Sprintf("/email/validate/%s", url.PathEscape(task.JobID)), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToHtml renders a validated email to HTML.
func (e *EmailEp) ToHtml(ctx context.Context, req EmailToHtmlRequest, opts *RequestOptions) (*EmailToHtmlResponse, error) {
	var out EmailToHtmlResponse
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, "/email/to/html", req, &task, to); err != nil {
		return nil, err
	}
	if err := e.t.pollResult(ctx, fmt.Sprintf("/email/to/html/%s", url.PathEscape(task.JobID)), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToPdf renders a validated email to PDF.
func (e *EmailEp) ToPdf(ctx context.Context, req EmailToPdfRequest, opts *RequestOptions) (*EmailToPdfResponse, error) {
	var out EmailToPdfResponse
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, "/email/to/pdf", req, &task, to); err != nil {
		return nil, err
	}
	if err := e.t.pollResult(ctx, fmt.Sprintf("/email/to/pdf/%s", url.PathEscape(task.JobID)), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}
