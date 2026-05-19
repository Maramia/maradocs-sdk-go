package maradocs

// ImgValidateRequest validates an uploaded image.
type ImgValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

type imgValidateInner struct {
	ClassName string    `json:"class_name"`
	ImgHandle *ImgHandle `json:"img_handle,omitempty"`
	Error     string    `json:"error,omitempty"`
	Virus     string    `json:"virus,omitempty"`
}

// ImgValidateResponse is the polled response for /img/validate.
type ImgValidateResponse struct {
	ClassName string           `json:"class_name"`
	Response  imgValidateInner `json:"response"`
}

// OkImg returns the image handle or an error matching the TypeScript okImg helper.
func OkImg(v ImgValidateResponse) (ImgHandle, error) {
	switch v.Response.ClassName {
	case "ImgValidateResponseOk":
		if v.Response.ImgHandle == nil {
			return ImgHandle{}, &ValidationError{Message: "missing img_handle"}
		}
		return *v.Response.ImgHandle, nil
	case "ImgValidateResponseError":
		return ImgHandle{}, &ValidationError{Message: v.Response.Error}
	case "ImgValidateResponseVirus":
		return ImgHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return ImgHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}

// ImgOrientationRequest detects orientation from text in the image.
type ImgOrientationRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgOrientationResponse contains the rotated image and detected angle.
type ImgOrientationResponse struct {
	RotatedImgHandle ImgHandle `json:"rotated_img_handle"`
	Orientation      int       `json:"orientation"`
}

// ImgRotateRequest rotates an image.
type ImgRotateRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
	Rotate    int       `json:"rotate"`
}

// ImgRotateResponse is the result of rotation.
type ImgRotateResponse struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgImproveContrastRequest improves contrast for document images.
type ImgImproveContrastRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgImproveContrastResponse is the enhanced image handle.
type ImgImproveContrastResponse struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgThumbnailRequest generates a thumbnail.
type ImgThumbnailRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
	MaxWidth  int       `json:"max_width"`
	MaxHeight int       `json:"max_height"`
}

// ImgThumbnailResponse is the thumbnail handle.
type ImgThumbnailResponse struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgToPngRequest converts to PNG.
type ImgToPngRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgToPngResponse is the PNG handle.
type ImgToPngResponse struct {
	PngHandle PngHandle `json:"png_handle"`
}

// ImgToJpegRequest converts to JPEG.
type ImgToJpegRequest struct {
	ImgHandle   ImgHandle `json:"img_handle"`
	Quality     *int      `json:"quality,omitempty"`
	Progressive *bool     `json:"progressive,omitempty"`
}

// ImgToJpegResponse is the JPEG handle.
type ImgToJpegResponse struct {
	JpegHandle JpegHandle `json:"jpeg_handle"`
}

// PdfImgColor is the color mode for PDF image embedding.
type PdfImgColor string

const (
	PdfImgColorOriginal   PdfImgColor = "original"
	PdfImgColorGrayscale  PdfImgColor = "grayscale"
	PdfImgColorMonochrome PdfImgColor = "monochrome"
)

// ImgToPdfOptions controls image-to-PDF conversion (subset of TS ImgToPdfOptionsBaseSchema).
type ImgToPdfOptions struct {
	ImgQuality        *int         `json:"img_quality,omitempty"`
	ImgColor          *PdfImgColor `json:"img_color,omitempty"`
	MaxSize           *PdfPageSize `json:"max_size,omitempty"`
	MaxDpi            *int         `json:"max_dpi,omitempty"`
	MinDpi            *int         `json:"min_dpi,omitempty"`
	EmbedInBlankPage  *struct {
		Size     PdfPageSize `json:"size"`
		Position string      `json:"position"`
	} `json:"embed_in_blank_page,omitempty"`
}

// ImgToPdfRequest converts an image to PDF.
type ImgToPdfRequest struct {
	ImgHandle ImgHandle        `json:"img_handle"`
	Options   *ImgToPdfOptions `json:"options,omitempty"`
}

// ImgToPdfResponse is the resulting PDF handle.
type ImgToPdfResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// ImgFindDocumentsRequest finds document regions.
type ImgFindDocumentsRequest struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgFindDocumentsQuadrilateral is one detected document region.
type ImgFindDocumentsQuadrilateral struct {
	Quadrilateral Quadrilateral `json:"quadrilateral"`
	Confidence      float64       `json:"confidence"`
}

// ImgFindDocumentsResponse lists detected documents.
type ImgFindDocumentsResponse struct {
	Documents []ImgFindDocumentsQuadrilateral `json:"documents"`
}

// ImgExtractQuadrilateralRequest extracts a perspective-corrected region.
type ImgExtractQuadrilateralRequest struct {
	ImgHandle     ImgHandle     `json:"img_handle"`
	Quadrilateral Quadrilateral `json:"quadrilateral"`
}

// ImgExtractQuadrilateralResponse is the extracted image handle.
type ImgExtractQuadrilateralResponse struct {
	ImgHandle ImgHandle `json:"img_handle"`
}

// ImgOcrToPdfRequest runs OCR and produces a searchable PDF.
type ImgOcrToPdfRequest struct {
	ImgHandle ImgHandle        `json:"img_handle"`
	Options   *ImgToPdfOptions `json:"options,omitempty"`
}

// ImgOcrToPdfResponse is the OCR'd PDF handle.
type ImgOcrToPdfResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}
