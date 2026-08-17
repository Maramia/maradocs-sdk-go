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

func boolPtr(v bool) *bool { return &v }

func requireProxyURL(proxyURL *string, kind string) (string, error) {
	if proxyURL == nil || *proxyURL == "" {
		return "", fmt.Errorf("%s proxy mint returned no proxy_url", kind)
	}
	return *proxyURL, nil
}

// CreateUpload mints a proxy-only upload capability URL for an unauthenticated third party.
// Always sets use_proxy; never returns post_url/post_header (SSE-C key material).
func (d *DataEp) CreateUpload(ctx context.Context, req DataUploadRequest) (*CreateUploadResult, error) {
	req.UseProxy = boolPtr(true)
	var resp DataUploadResponse
	if err := d.t.postJSON(ctx, "/data/upload", &req, &resp, nil); err != nil {
		return nil, err
	}
	proxyURL, err := requireProxyURL(resp.ProxyURL, "upload")
	if err != nil {
		return nil, err
	}
	return &CreateUploadResult{
		ProxyURL:              proxyURL,
		UnvalidatedFileHandle: resp.UnvalidatedFileHandle,
	}, nil
}

// Upload requests a presigned POST and uploads the bytes to object storage (first-party, no proxy).
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

func (d *DataEp) createDownloadFromMeta(proxyURL *string) (*CreateDownloadProxyResult, error) {
	url, err := requireProxyURL(proxyURL, "download")
	if err != nil {
		return nil, err
	}
	return &CreateDownloadProxyResult{ProxyURL: url}, nil
}

// CreateDownloadPdf mints a proxy-only PDF download capability URL.
func (d *DataEp) CreateDownloadPdf(ctx context.Context, req DataDownloadPdfRequest) (*CreateDownloadProxyResult, error) {
	req.UseProxy = boolPtr(true)
	var meta DataDownloadPdfResponse
	if err := d.t.postJSON(ctx, "/data/download/pdf", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.createDownloadFromMeta(meta.ProxyURL)
}

// DownloadPdf returns raw PDF bytes (first-party SSE-C).
func (d *DataEp) DownloadPdf(ctx context.Context, req DataDownloadPdfRequest, onProgress func(float64)) ([]byte, error) {
	req.UseProxy = nil
	var meta DataDownloadPdfResponse
	if err := d.t.postJSON(ctx, "/data/download/pdf", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// CreateDownloadJpeg mints a proxy-only JPEG download capability URL.
func (d *DataEp) CreateDownloadJpeg(ctx context.Context, req DataDownloadJpegRequest) (*CreateDownloadProxyResult, error) {
	req.UseProxy = boolPtr(true)
	var meta DataDownloadJpegResponse
	if err := d.t.postJSON(ctx, "/data/download/jpeg", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.createDownloadFromMeta(meta.ProxyURL)
}

// DownloadJpeg returns raw JPEG bytes (first-party SSE-C).
func (d *DataEp) DownloadJpeg(ctx context.Context, req DataDownloadJpegRequest, onProgress func(float64)) ([]byte, error) {
	req.UseProxy = nil
	var meta DataDownloadJpegResponse
	if err := d.t.postJSON(ctx, "/data/download/jpeg", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// CreateDownloadPng mints a proxy-only PNG download capability URL.
func (d *DataEp) CreateDownloadPng(ctx context.Context, req DataDownloadPngRequest) (*CreateDownloadProxyResult, error) {
	req.UseProxy = boolPtr(true)
	var meta DataDownloadPngResponse
	if err := d.t.postJSON(ctx, "/data/download/png", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.createDownloadFromMeta(meta.ProxyURL)
}

// DownloadPng returns raw PNG bytes (first-party SSE-C).
func (d *DataEp) DownloadPng(ctx context.Context, req DataDownloadPngRequest, onProgress func(float64)) ([]byte, error) {
	req.UseProxy = nil
	var meta DataDownloadPngResponse
	if err := d.t.postJSON(ctx, "/data/download/png", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// CreateDownloadOdt mints a proxy-only ODT download capability URL.
func (d *DataEp) CreateDownloadOdt(ctx context.Context, req DataDownloadOdtRequest) (*CreateDownloadProxyResult, error) {
	req.UseProxy = boolPtr(true)
	var meta DataDownloadOdtResponse
	if err := d.t.postJSON(ctx, "/data/download/odt", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.createDownloadFromMeta(meta.ProxyURL)
}

// DownloadOdt returns raw ODT bytes (first-party SSE-C).
func (d *DataEp) DownloadOdt(ctx context.Context, req DataDownloadOdtRequest, onProgress func(float64)) ([]byte, error) {
	req.UseProxy = nil
	var meta DataDownloadOdtResponse
	if err := d.t.postJSON(ctx, "/data/download/odt", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// CreateDownloadUnvalidated mints a proxy-only unvalidated-file download capability URL.
func (d *DataEp) CreateDownloadUnvalidated(ctx context.Context, req DataDownloadUnvalidatedRequest) (*CreateDownloadProxyResult, error) {
	req.UseProxy = boolPtr(true)
	var meta DataDownloadUnvalidatedResponse
	if err := d.t.postJSON(ctx, "/data/download/unvalidated", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.createDownloadFromMeta(meta.ProxyURL)
}

// DownloadUnvalidated downloads an unvalidated file by first-party presigned URL.
func (d *DataEp) DownloadUnvalidated(ctx context.Context, req DataDownloadUnvalidatedRequest, onProgress func(float64)) ([]byte, error) {
	req.UseProxy = nil
	var meta DataDownloadUnvalidatedResponse
	if err := d.t.postJSON(ctx, "/data/download/unvalidated", req, &meta, nil); err != nil {
		return nil, err
	}
	return d.downloadGET(ctx, meta.URL, meta.Headers, onProgress)
}

// DownloadMp4 runs async export then downloads MP4 bytes.
func (d *DataEp) DownloadMp4(ctx context.Context, req DataDownloadMp4Request, onProgress func(float64), opts *RequestOptions) ([]byte, error) {
	req.UseProxy = nil
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
	req.UseProxy = nil
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
	req.UseProxy = nil
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
	req.UseProxy = nil
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
