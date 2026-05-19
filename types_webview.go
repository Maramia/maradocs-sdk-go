package maradocs

// WebviewOpenRequest opens the interactive webview.
type WebviewOpenRequest struct {
	Language          *string       `json:"language,omitempty"`
	AllowUserUpload   *bool         `json:"allow_user_upload,omitempty"`
	AllowUserDownload *bool         `json:"allow_user_download,omitempty"`
	WithFiles         []FileHandle  `json:"with_files,omitempty"`
}

// WebviewOpenResponse returns a browser URL.
type WebviewOpenResponse struct {
	URL string `json:"url"`
}

// WebviewAddFileRequest adds a file to the webview session.
type WebviewAddFileRequest struct {
	FileHandle FileHandle `json:"file_handle"`
}

// WebviewAddFileResponse is an empty object for backwards compatibility.
type WebviewAddFileResponse struct{}

// WebviewStatusResponse reports whether the webview is open.
type WebviewStatusResponse struct {
	Status string `json:"status"`
}

// WebviewResultsResponse returns files produced in a webview session.
type WebviewResultsResponse struct {
	UploadedPdfs  []PdfHandle `json:"uploaded_pdfs"`
	UploadedImgs  []ImgHandle `json:"uploaded_imgs"`
	ProcessedPdfs []PdfHandle `json:"processed_pdfs"`
}
