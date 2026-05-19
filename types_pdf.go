package maradocs

// PdfValidateRequest validates an uploaded PDF.
type PdfValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
	Password              *string               `json:"password,omitempty"`
}

type pdfValidateInner struct {
	ClassName string     `json:"class_name"`
	PdfHandle *PdfHandle `json:"pdf_handle,omitempty"`
	Error     string     `json:"error,omitempty"`
	Virus     string     `json:"virus,omitempty"`
}

// PdfValidateResponse is the polled response for /pdf/validate.
type PdfValidateResponse struct {
	ClassName string           `json:"class_name"`
	Response  pdfValidateInner `json:"response"`
}

// OkPdf extracts the PDF handle from a validation response.
func OkPdf(v PdfValidateResponse) (PdfHandle, error) {
	switch v.Response.ClassName {
	case "PdfValidateResponseOk":
		if v.Response.PdfHandle == nil {
			return PdfHandle{}, &ValidationError{Message: "missing pdf_handle"}
		}
		return *v.Response.PdfHandle, nil
	case "PdfValidateResponseError":
		return PdfHandle{}, &ValidationError{Message: v.Response.Error}
	case "PdfValidateResponseVirus":
		return PdfHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return PdfHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}

// PdfComposePdfPage selects a page with optional rotation.
type PdfComposePdfPage struct {
	PageNumber int  `json:"page_number"`
	Rotation   *int `json:"rotation,omitempty"`
}

// PdfComposePdf is one PDF source in a composition.
// Pages nil omits the field (all pages). Use a non-nil pointer to send an explicit
// page list, including an empty slice for zero pages (plain [] with omitempty would not serialize []).
type PdfComposePdf struct {
	PdfHandle PdfHandle            `json:"pdf_handle"`
	Pages     *[]PdfComposePdfPage `json:"pages,omitempty"`
}

// PdfComposeRequest composes PDFs.
type PdfComposeRequest struct {
	Pdfs []PdfComposePdf `json:"pdfs"`
}

// PdfComposeResponse is the composed PDF handle.
type PdfComposeResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// PdfRotateRequest rotates selected pages.
type PdfRotateRequest struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
	Rotate    [][]int   `json:"rotate"`
}

// PdfRotateResponse is the rotated PDF handle.
type PdfRotateResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// PdfOptimizeRequest optimizes a PDF.
type PdfOptimizeRequest struct {
	PdfHandle    PdfHandle    `json:"pdf_handle"`
	ImageDpi     *int         `json:"image_dpi,omitempty"`
	ImageQuality *int         `json:"image_quality,omitempty"`
	ImageColor   *PdfImgColor `json:"image_color,omitempty"`
}

// PdfOptimizeResponse is the optimized PDF handle.
type PdfOptimizeResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// PdfToImgRequest renders PDF pages as images.
type PdfToImgRequest struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
	Pages     []int     `json:"pages,omitempty"`
	Dpi       *int      `json:"dpi,omitempty"`
}

// PdfToImgResponse lists rendered page images.
type PdfToImgResponse struct {
	ImgHandles []ImgHandle `json:"img_handles"`
}

// PdfOrientationRequest detects per-page orientation.
type PdfOrientationRequest struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// PdfOrientationResponse contains per-page angles and rotated PDF.
type PdfOrientationResponse struct {
	Orientations     []int     `json:"orientations"`
	RotatedPdfHandle PdfHandle `json:"rotated_pdf_handle"`
}

// PdfOcrToPdfRequest OCRs a PDF to searchable PDF.
type PdfOcrToPdfRequest struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}

// PdfOcrToPdfResponse is the OCR'd PDF handle.
type PdfOcrToPdfResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}
