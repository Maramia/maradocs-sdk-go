package maradocs

import (
	"context"
	"fmt"
	"net/url"
)

// ImgEp implements /img/* endpoints.
type ImgEp struct {
	t *transport
}

func newImgEp(t *transport) *ImgEp {
	return &ImgEp{t: t}
}

func (e *ImgEp) postPoll(ctx context.Context, postPath string, req any, pollFmt string, out any, opts *RequestOptions) error {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := e.t.postJSON(ctx, postPath, req, &task, to); err != nil {
		return err
	}
	pollPath := fmt.Sprintf(pollFmt, url.PathEscape(task.JobID))
	return e.t.pollResult(ctx, pollPath, out, to)
}

// Validate validates an uploaded image.
func (e *ImgEp) Validate(ctx context.Context, req ImgValidateRequest, opts *RequestOptions) (*ImgValidateResponse, error) {
	var out ImgValidateResponse
	if err := e.postPoll(ctx, "/img/validate", req, "/img/validate/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Thumbnail generates a thumbnail image.
func (e *ImgEp) Thumbnail(ctx context.Context, req ImgThumbnailRequest, opts *RequestOptions) (*ImgThumbnailResponse, error) {
	var out ImgThumbnailResponse
	if err := e.postPoll(ctx, "/img/thumbnail", req, "/img/thumbnail/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// FindDocuments detects document regions.
func (e *ImgEp) FindDocuments(ctx context.Context, req ImgFindDocumentsRequest, opts *RequestOptions) (*ImgFindDocumentsResponse, error) {
	var out ImgFindDocumentsResponse
	if err := e.postPoll(ctx, "/img/find/documents", req, "/img/find/documents/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Orientation detects orientation from depicted text.
func (e *ImgEp) Orientation(ctx context.Context, req ImgOrientationRequest, opts *RequestOptions) (*ImgOrientationResponse, error) {
	var out ImgOrientationResponse
	if err := e.postPoll(ctx, "/img/orientation", req, "/img/orientation/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExtractQuadrilateral extracts a perspective-corrected region.
func (e *ImgEp) ExtractQuadrilateral(ctx context.Context, req ImgExtractQuadrilateralRequest, opts *RequestOptions) (*ImgExtractQuadrilateralResponse, error) {
	var out ImgExtractQuadrilateralResponse
	if err := e.postPoll(ctx, "/img/extract/quadrilateral", req, "/img/extract/quadrilateral/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rotate rotates an image by a discrete angle.
func (e *ImgEp) Rotate(ctx context.Context, req ImgRotateRequest, opts *RequestOptions) (*ImgRotateResponse, error) {
	var out ImgRotateResponse
	if err := e.postPoll(ctx, "/img/rotate", req, "/img/rotate/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ImproveContrast improves lighting for document images.
func (e *ImgEp) ImproveContrast(ctx context.Context, req ImgImproveContrastRequest, opts *RequestOptions) (*ImgImproveContrastResponse, error) {
	var out ImgImproveContrastResponse
	if err := e.postPoll(ctx, "/img/improve-contrast", req, "/img/improve-contrast/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToJpeg converts an image to JPEG.
func (e *ImgEp) ToJpeg(ctx context.Context, req ImgToJpegRequest, opts *RequestOptions) (*ImgToJpegResponse, error) {
	var out ImgToJpegResponse
	if err := e.postPoll(ctx, "/img/to/jpeg", req, "/img/to/jpeg/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToPng converts an image to PNG.
func (e *ImgEp) ToPng(ctx context.Context, req ImgToPngRequest, opts *RequestOptions) (*ImgToPngResponse, error) {
	var out ImgToPngResponse
	if err := e.postPoll(ctx, "/img/to/png", req, "/img/to/png/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// ToPdf converts an image to PDF (without OCR text layer).
func (e *ImgEp) ToPdf(ctx context.Context, req ImgToPdfRequest, opts *RequestOptions) (*ImgToPdfResponse, error) {
	var out ImgToPdfResponse
	if err := e.postPoll(ctx, "/img/to/pdf", req, "/img/to/pdf/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

// OcrToPdf runs OCR and returns a searchable PDF handle.
func (e *ImgEp) OcrToPdf(ctx context.Context, req ImgOcrToPdfRequest, opts *RequestOptions) (*ImgOcrToPdfResponse, error) {
	var out ImgOcrToPdfResponse
	if err := e.postPoll(ctx, "/img/ocr/to/pdf", req, "/img/ocr/to/pdf/%s", &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}
