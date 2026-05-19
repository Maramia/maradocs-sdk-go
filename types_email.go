package maradocs

import "encoding/json"

// NamedAddress is an email address with optional display name.
type NamedAddress struct {
	Name    *string `json:"name"`
	Address string  `json:"address"`
}

// EmailAttachment describes one attachment on a parsed email.
type EmailAttachment struct {
	ClassName             string                `json:"class_name"`
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
	MediaType             DataMediaTypeResponse `json:"media_type"`
	Validated             json.RawMessage       `json:"validated"`
	Name                  *string               `json:"name"`
	ContentID             *string               `json:"content_id"`
	ContentDisposition    *string               `json:"content_disposition"`
}

// EmailHandle is a validated email parse result.
type EmailHandle struct {
	ClassName   string                 `json:"class_name"`
	FileHandle  FileHandle             `json:"file_handle"`
	SignedHash  string                 `json:"signed_hash"`
	FromAddr    []NamedAddress         `json:"from_addr"`
	ToAddr      []NamedAddress         `json:"to_addr"`
	CcAddr      []NamedAddress         `json:"cc_addr"`
	BccAddr     []NamedAddress         `json:"bcc_addr"`
	Date        *string                `json:"date"`
	Subject     *string                `json:"subject"`
	TextBody    *UnvalidatedFileHandle `json:"text_body"`
	HtmlBody    *UnvalidatedFileHandle `json:"html_body"`
	Attachments []EmailAttachment      `json:"attachments"`
}

// EmailValidateRequest validates an uploaded email file.
type EmailValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

type emailValidateInner struct {
	ClassName   string        `json:"class_name"`
	EmailHandle *EmailHandle  `json:"email_handle,omitempty"`
	Error       string        `json:"error,omitempty"`
	Virus       string        `json:"virus,omitempty"`
}

// EmailValidateResponse is the polled response for /email/validate.
type EmailValidateResponse struct {
	ClassName string             `json:"class_name"`
	Response  emailValidateInner `json:"response"`
}

// OkEmail extracts the email handle from a validation response.
func OkEmail(v EmailValidateResponse) (EmailHandle, error) {
	switch v.Response.ClassName {
	case "EmailValidateResponseOk":
		if v.Response.EmailHandle == nil {
			return EmailHandle{}, &ValidationError{Message: "missing email_handle"}
		}
		return *v.Response.EmailHandle, nil
	case "EmailValidateResponseError":
		return EmailHandle{}, &ValidationError{Message: v.Response.Error}
	case "EmailValidateResponseVirus":
		return EmailHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return EmailHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}

// EmailTemplateLabels are custom labels for HTML rendering.
type EmailTemplateLabels struct {
	FromLabel string `json:"from_label"`
	ToLabel   string `json:"to_label"`
	CcLabel   string `json:"cc_label"`
	DateLabel string `json:"date_label"`
}

// EmailToHtmlOptions configures email HTML rendering.
type EmailToHtmlOptions struct {
	Locale *string              `json:"locale,omitempty"`
	Labels *EmailTemplateLabels `json:"labels,omitempty"`
}

// EmailToHtmlRequest renders email to HTML.
type EmailToHtmlRequest struct {
	EmailHandle EmailHandle         `json:"email_handle"`
	Options     *EmailToHtmlOptions `json:"options,omitempty"`
}

// EmailToHtmlResponse is the rendered HTML file handle.
type EmailToHtmlResponse struct {
	HtmlHandle HtmlHandle `json:"html_handle"`
}

// EmailToPdfRequest renders email to PDF.
type EmailToPdfRequest struct {
	EmailHandle EmailHandle         `json:"email_handle"`
	OptionsHTML *EmailToHtmlOptions `json:"options_html,omitempty"`
}

// EmailToPdfResponse is the rendered PDF handle.
type EmailToPdfResponse struct {
	PdfHandle PdfHandle `json:"pdf_handle"`
}
