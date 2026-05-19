package maradocs

import (
	"context"
	"fmt"
	"net/url"
)

// PdfEp implements /pdf/* endpoints.
type PdfEp struct {
	t *transport
}

func newPdfEp(t *transport) *PdfEp {
	return &PdfEp{t: t}
}

func (e *PdfEp) postPoll(ctx context.Context, postPath string, req any, pollFmt string, out any, opts *RequestOptions) error {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, postPath, req, &task, to); err != nil {
		return err
	}
	return e.t.pollResult(ctx, fmt.Sprintf(pollFmt, url.PathEscape(task.JobID)), out, to)
}

// Validate validates an uploaded PDF.
func (e *PdfEp) Validate(ctx context.Context, req PdfValidateRequest, opts *RequestOptions) (*PdfValidateResponse, error) {
	var out PdfValidateResponse
	if err := e.postPoll(ctx, "/pdf/validate", req, "/pdf/validate/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Compose composes PDFs from existing PDFs/pages.
func (e *PdfEp) Compose(ctx context.Context, req PdfComposeRequest, opts *RequestOptions) (*PdfComposeResponse, error) {
	var out PdfComposeResponse
	if err := e.postPoll(ctx, "/pdf/compose", req, "/pdf/compose/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Optimize optimizes a PDF.
func (e *PdfEp) Optimize(ctx context.Context, req PdfOptimizeRequest, opts *RequestOptions) (*PdfOptimizeResponse, error) {
	var out PdfOptimizeResponse
	if err := e.postPoll(ctx, "/pdf/optimize", req, "/pdf/optimize/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rotate rotates selected PDF pages.
func (e *PdfEp) Rotate(ctx context.Context, req PdfRotateRequest, opts *RequestOptions) (*PdfRotateResponse, error) {
	var out PdfRotateResponse
	if err := e.postPoll(ctx, "/pdf/rotate", req, "/pdf/rotate/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToImg renders PDF pages as images.
func (e *PdfEp) ToImg(ctx context.Context, req PdfToImgRequest, opts *RequestOptions) (*PdfToImgResponse, error) {
	var out PdfToImgResponse
	if err := e.postPoll(ctx, "/pdf/to/img", req, "/pdf/to/img/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Orientation detects and corrects PDF page orientation.
func (e *PdfEp) Orientation(ctx context.Context, req PdfOrientationRequest, opts *RequestOptions) (*PdfOrientationResponse, error) {
	var out PdfOrientationResponse
	if err := e.postPoll(ctx, "/pdf/orientation", req, "/pdf/orientation/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// OcrToPdf OCRs a PDF to a searchable PDF.
func (e *PdfEp) OcrToPdf(ctx context.Context, req PdfOcrToPdfRequest, opts *RequestOptions) (*PdfOcrToPdfResponse, error) {
	var out PdfOcrToPdfResponse
	if err := e.postPoll(ctx, "/pdf/ocr/pdf", req, "/pdf/ocr/pdf/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}
