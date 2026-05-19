package maradocs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// DataEp implements /data/* workspace endpoints.
type DataEp struct {
	t *transport
}

func newDataEp(t *transport) *DataEp {
	return &DataEp{t: t}
}

func defaultProgress(onProgress func(float64)) func(float64) {
	if onProgress == nil {
		return func(float64) {}
	}
	return onProgress
}

// Upload requests a presigned POST and uploads the bytes to object storage.
func (d *DataEp) Upload(ctx context.Context, fileName string, size int64, r io.Reader, onProgress func(float64)) (*DataUploadResponse, error) {
	onProgress = defaultProgress(onProgress)
	req := DataUploadRequest{Size: size}
	if fileName != "" {
		req.Name = &fileName
	}
	var resp DataUploadResponse
	if err := d.t.postJSON(ctx, "/data/upload", &req, &resp, nil); err != nil {
		return nil, err
	}
	uploadBody := &bytes.Buffer{}
	mw := multipart.NewWriter(uploadBody)
	for k, v := range resp.PostHeader {
		if err := mw.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	fn := fileName
	if fn == "" {
		fn = "upload"
	}
	fw, err := mw.CreateFormFile("file", fn)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, r); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, resp.PostURL, bytes.NewReader(uploadBody.Bytes()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", mw.FormDataContentType())
	httpResp, err := d.t.doHTTP(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		bb, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("S3 upload failed: %s %s", httpResp.Status, truncate(bb, 256))
	}
	_, _ = io.Copy(io.Discard, httpResp.Body)
	onProgress(100)
	return &resp, nil
}

func (d *DataEp) downloadGET(ctx context.Context, url string, headers map[string]string, onProgress func(float64)) ([]byte, error) {
	onProgress = defaultProgress(onProgress)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := d.t.doHTTP(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bb, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed: %s %s", resp.Status, truncate(bb, 256))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	onProgress(100)
	return b, nil
}

// VirusScan runs virus scan task + poll.
func (d *DataEp) VirusScan(ctx context.Context, req VirusScanRequest, opts *RequestOptions) (*VirusScanResponse, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/virus_scan", req, &task, to); err != nil {
		return nil, err
	}
	var out VirusScanResponse
	path := fmt.Sprintf("/data/virus_scan/%s", task.JobID)
	if err := d.t.pollResult(ctx, path, &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}

// MimeType runs MIME detection task + poll.
func (d *DataEp) MimeType(ctx context.Context, req DataMediaTypeRequest, opts *RequestOptions) (*DataMediaTypeResponse, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/mime_type", req, &task, to); err != nil {
		return nil, err
	}
	var out DataMediaTypeResponse
	if err := d.t.pollResult(ctx, fmt.Sprintf("/data/mime_type/%s", task.JobID), &out, to); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadPdf returns raw PDF bytes.
func (d *DataEp) DownloadPdf(ctx context.Context, req DataDownloadPdfRequest, onProgress func(float64)) ([]byte, error) {
	var meta DataDownloadPdfResponse
	if err := d.t.postJSON(ctx, "/data/download/pdf", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadJpeg returns raw JPEG bytes.
func (d *DataEp) DownloadJpeg(ctx context.Context, req DataDownloadJpegRequest, onProgress func(float64)) ([]byte, error) {
	var meta DataDownloadJpegResponse
	if err := d.t.postJSON(ctx, "/data/download/jpeg", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadPng returns raw PNG bytes.
func (d *DataEp) DownloadPng(ctx context.Context, req DataDownloadPngRequest, onProgress func(float64)) ([]byte, error) {
	var meta DataDownloadPngResponse
	if err := d.t.postJSON(ctx, "/data/download/png", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadOdt returns raw ODT bytes.
func (d *DataEp) DownloadOdt(ctx context.Context, req DataDownloadOdtRequest, onProgress func(float64)) ([]byte, error) {
	var meta DataDownloadOdtResponse
	if err := d.t.postJSON(ctx, "/data/download/odt", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadUnvalidated downloads an unvalidated file by presigned URL.
func (d *DataEp) DownloadUnvalidated(ctx context.Context, req DataDownloadUnvalidatedRequest, onProgress func(float64)) ([]byte, error) {
	var meta DataDownloadUnvalidatedResponse
	if err := d.t.postJSON(ctx, "/data/download/unvalidated", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadMp4 runs async export then downloads MP4 bytes.
func (d *DataEp) DownloadMp4(ctx context.Context, req DataDownloadMp4Request, onProgress func(float64), opts *RequestOptions) ([]byte, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/download/mp4", req, &task, to); err != nil {
		return nil, err
	}
	var meta DataDownloadMp4Response
	if err := d.t.pollResult(ctx, fmt.Sprintf("/data/download/mp4/%s", task.JobID), &meta, to); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadMp3 runs async export then downloads MP3 bytes.
func (d *DataEp) DownloadMp3(ctx context.Context, req DataDownloadMp3Request, onProgress func(float64), opts *RequestOptions) ([]byte, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/download/mp3", req, &task, to); err != nil {
		return nil, err
	}
	var meta DataDownloadMp3Response
	if err := d.t.pollResult(ctx, fmt.Sprintf("/data/download/mp3/%s", task.JobID), &meta, to); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadWav runs async export then downloads WAV bytes.
func (d *DataEp) DownloadWav(ctx context.Context, req DataDownloadWavRequest, onProgress func(float64), opts *RequestOptions) ([]byte, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/download/wav", req, &task, to); err != nil {
		return nil, err
	}
	var meta DataDownloadWavResponse
	if err := d.t.pollResult(ctx, fmt.Sprintf("/data/download/wav/%s", task.JobID), &meta, to); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadFlac runs async export then downloads FLAC bytes.
func (d *DataEp) DownloadFlac(ctx context.Context, req DataDownloadFlacRequest, onProgress func(float64), opts *RequestOptions) ([]byte, error) {
	to := optTimeoutMs(opts)
	var task TaskCreatedResponse
	if err := d.t.postJSON(ctx, "/data/download/flac", req, &task, to); err != nil {
		return nil, err
	}
	var meta DataDownloadFlacResponse
	if err := d.t.pollResult(ctx, fmt.Sprintf("/data/download/flac/%s", task.JobID), &meta, to); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}
