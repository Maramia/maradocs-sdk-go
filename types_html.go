package maradocs

// HtmlValidateRequest validates uploaded HTML.
type HtmlValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

type htmlValidateInner struct {
	ClassName  string       `json:"class_name"`
	HtmlHandle *HtmlHandle  `json:"html_handle,omitempty"`
	Error      string       `json:"error,omitempty"`
	Virus      string       `json:"virus,omitempty"`
}

// HtmlValidateResponse is the polled response for /html/validate.
type HtmlValidateResponse struct {
	ClassName string             `json:"class_name"`
	Response  htmlValidateInner  `json:"response"`
}

// OkHtml extracts the HTML handle from a validation response.
func OkHtml(v HtmlValidateResponse) (HtmlHandle, error) {
	switch v.Response.ClassName {
	case "HtmlValidateResponseOk":
		if v.Response.HtmlHandle == nil {
			return HtmlHandle{}, &ValidationError{Message: "missing html_handle"}
		}
		return *v.Response.HtmlHandle, nil
	case "HtmlValidateResponseError":
		return HtmlHandle{}, &ValidationError{Message: v.Response.Error}
	case "HtmlValidateResponseVirus":
		return HtmlHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return HtmlHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}

// HtmlToPdfRequest converts HTML to PDF.
type HtmlToPdfRequest struct {
	HtmlHandle HtmlHandle `json:"html_handle"`
}

// HtmlToPdfResponse is the converted PDF handle.
type HtmlToPdfResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}
