package maradocs

// DataUploadRequest starts a workspace file upload.
type DataUploadRequest struct {
	Name *string `json:"name,omitempty"`
	Size int64   `json:"size"`
}

// DataUploadResponse contains presigned POST data for S3 and the unvalidated handle.
type DataUploadResponse struct {
	PostURL               string                `json:"post_url"`
	PostHeader            map[string]string     `json:"post_header"`
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

// DataDownloadPdfRequest requests a presigned PDF download URL.
type DataDownloadPdfRequest struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
	ExpiresIn *int      `json:"expires_in,omitempty"`
}

// DataDownloadPdfResponse contains presigned GET URL and headers.
type DataDownloadPdfResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadJpegRequest requests a JPEG download URL.
type DataDownloadJpegRequest struct {
	JpegHandle JpegHandle `json:"jpeg_handle"`
	ExpiresIn  *int       `json:"expires_in,omitempty"`
}

// DataDownloadJpegResponse contains presigned GET URL and headers.
type DataDownloadJpegResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadPngRequest requests a PNG download URL.
type DataDownloadPngRequest struct {
	PngHandle PngHandle `json:"png_handle"`
	ExpiresIn *int      `json:"expires_in,omitempty"`
}

// DataDownloadPngResponse contains presigned GET URL and headers.
type DataDownloadPngResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadOdtRequest requests an ODT download URL.
type DataDownloadOdtRequest struct {
	OdtHandle OdtHandle `json:"odt_handle"`
	ExpiresIn *int      `json:"expires_in,omitempty"`
}

// DataDownloadOdtResponse contains presigned GET URL and headers.
type DataDownloadOdtResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadUnvalidatedRequest downloads an unvalidated file.
type DataDownloadUnvalidatedRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
	ExpiresIn             *int                  `json:"expires_in,omitempty"`
}

// DataDownloadUnvalidatedResponse contains presigned GET URL and headers.
type DataDownloadUnvalidatedResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataMediaTypeRequest detects MIME type of an unvalidated file.
type DataMediaTypeRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

// DataMediaTypeResponse is the MIME detection result.
type DataMediaTypeResponse struct {
	MediaType  string  `json:"media_type"`
	Confidence float64 `json:"confidence"`
}

// VirusScanRequest scans a file for malware.
type VirusScanRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

// VirusScanResponse is the virus scan result.
type VirusScanResponse struct {
	VirusFound bool    `json:"virus_found"`
	VirusInfo  *string `json:"virus_info"`
}

// DataDownloadMp4Request requests MP4 export/transcode.
type DataDownloadMp4Request struct {
	VideoHandle        VideoHandle `json:"video_handle"`
	ConstantFrameRate  *string     `json:"constant_frame_rate,omitempty"`
	AudioCodec         *string     `json:"audio_codec,omitempty"`
	AudioBitrate       *string     `json:"audio_bitrate,omitempty"`
	Format             *string     `json:"format,omitempty"`
	ExpiresIn          *int        `json:"expires_in,omitempty"`
}

// DataDownloadMp4Response contains presigned GET URL and headers.
type DataDownloadMp4Response struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadMp3Request requests MP3 export.
type DataDownloadMp3Request struct {
	AudioHandle AudioHandle `json:"audio_handle"`
	Bitrate     *string     `json:"bitrate,omitempty"`
	Format      *string     `json:"format,omitempty"`
	ExpiresIn   *int        `json:"expires_in,omitempty"`
}

// DataDownloadMp3Response contains presigned GET URL and headers.
type DataDownloadMp3Response struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadWavRequest requests WAV export.
type DataDownloadWavRequest struct {
	AudioHandle AudioHandle `json:"audio_handle"`
	BitDepth    *string     `json:"bit_depth,omitempty"`
	Format      *string     `json:"format,omitempty"`
	ExpiresIn   *int        `json:"expires_in,omitempty"`
}

// DataDownloadWavResponse contains presigned GET URL and headers.
type DataDownloadWavResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// DataDownloadFlacRequest requests FLAC export.
type DataDownloadFlacRequest struct {
	AudioHandle       AudioHandle `json:"audio_handle"`
	CompressionLevel  *int        `json:"compression_level,omitempty"`
	Format            *string     `json:"format,omitempty"`
	ExpiresIn         *int        `json:"expires_in,omitempty"`
}

// DataDownloadFlacResponse contains presigned GET URL and headers.
type DataDownloadFlacResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
