package maradocs

// FileHandle is a reference to a file in workspace object storage.
type FileHandle struct {
	Reference string `json:"reference"`
	Size      int64  `json:"size"`
	ExpiresAt string `json:"expires_at"`
}

// UnvalidatedFileHandle is a file that has not been validated yet.
type UnvalidatedFileHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
}

// ImgHandle is a validated image handle.
type ImgHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
}

// PdfPageSize is PDF page size in millimeters.
type PdfPageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// PdfHandle is a validated PDF handle.
type PdfHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Pages      int        `json:"pages"`
}

// PngHandle is a PNG image handle.
type PngHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
}

// JpegHandle is a JPEG image handle.
type JpegHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
}

// OdtHandle is an ODT document handle.
type OdtHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Pages      int        `json:"pages"`
}

// HtmlHandle is a validated HTML file handle.
type HtmlHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
}

// AudioMetadata describes source audio stream metadata.
type AudioMetadata struct {
	SampleRateHz    int    `json:"sample_rate_hz"`
	Channels        int    `json:"channels"`
	BitDepth        int    `json:"bit_depth"`
	Codec           string `json:"codec"`
	ContainerFormat string `json:"container_format"`
	BitRateKbps     int    `json:"bit_rate_kbps"`
}

// VideoMetadata describes source video stream metadata.
type VideoMetadata struct {
	FrameRate        float64 `json:"frame_rate"`
	Codec            string  `json:"codec"`
	ContainerFormat  string  `json:"container_format"`
	BitRateKbps      *int    `json:"bit_rate_kbps"`
	PixelFormat      *string `json:"pixel_format"`
}

// VideoHandle is a validated video handle.
type VideoHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	DurationMs int       `json:"duration_ms"`
}

// AudioHandle is a validated audio handle.
type AudioHandle struct {
	SignedHash string     `json:"signed_hash"`
	FileHandle FileHandle `json:"file_handle"`
	DurationMs int        `json:"duration_ms"`
}

// TaskCreatedResponse is returned when an async job is created.
type TaskCreatedResponse struct {
	JobID       string `json:"job_id"`
	TokensSpent int    `json:"tokens_spent"`
}

// RelativePosition is a normalized point in an image (0..1).
type RelativePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Quadrilateral defines document corners in relative coordinates.
type Quadrilateral struct {
	TopLeft     RelativePosition `json:"top_left"`
	TopRight    RelativePosition `json:"top_right"`
	BottomRight RelativePosition `json:"bottom_right"`
	BottomLeft  RelativePosition `json:"bottom_left"`
}
