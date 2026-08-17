# MaraDocs Go SDK

[MaraDocs](https://maradocs.io) turns photos and scans into clean, usable documents: PDF conversion, virus checks, OCR, compression, and more. The MaraDocs API exposes those capabilities for your own automation.

This is the official Go client for the [MaraDocs API](https://api.maradocs.io). It mirrors the public TypeScript SDK (`maradocs-sdk-ts`): workspace and account clients, async task polling, validation helpers, and high-level `Flow` pipelines.

## Installation

```bash
go get github.com/maramia/maradocs-sdk-go@v0.2.0
```

Requires Go 1.22+.

## Documentation

Full API reference: [api.maradocs.io](https://api.maradocs.io)

## Quick start

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/maramia/maradocs-sdk-go"
)

func main() {
	ctx := context.Background()

	// Server: create a workspace and hand workspace_secret to the client
	srv := maradocs.NewMaraDocsServer(maradocs.MaraDocsServerOptions{
		SecretKey: os.Getenv("MARADOCS_SECRET_KEY"),
	})
	ws, err := srv.Workspace.Create(ctx, maradocs.WorkspaceCreateRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Client: OCR an image and combine with another PDF
	client, err := maradocs.NewMaraDocsClient(maradocs.MaraDocsClientOptions{
		WorkspaceSecret: ws.WorkspaceSecret,
	})
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("scan.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	st, _ := f.Stat()

	imgPdf, err := client.Flow.OcrImg(ctx, "scan.jpg", st.Size(), f, maradocs.OcrImgOptions{})
	if err != nil {
		log.Fatal(err)
	}

	combined, err := client.Pdf.Compose(ctx, maradocs.PdfComposeRequest{
		Pdfs: []maradocs.PdfComposePdf{{PdfHandle: imgPdf}},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	pdfBytes, err := client.Data.DownloadPdf(ctx, maradocs.DataDownloadPdfRequest{
		PdfHandle: combined.PdfHandle,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = pdfBytes
}
```

Every API call takes a `context.Context`. Set `APIURLWithVersion` on client options to override the default `https://api.maradocs.io/v1`. After `NewMaraDocsClient`, `client.Info` holds decoded workspace metadata (account id, workspace id, encryption key).

## Clients

| Client | Use case | Authentication |
|--------|----------|------------------|
| `MaraDocsClient` | Document processing in a workspace | Workspace secret |
| `MaraDocsServer` | Account, workspace, and webview management | Account secret key |
| `NewHealthcheckClient` | Unauthenticated `GET /healthcheck/ping` | None |

Optional `TimeoutMs` on client constructors sets a default timeout (milliseconds) for HTTP requests. Poll-based methods also accept `*maradocs.RequestOptions{ Timeout: &ms }` per call.

Downloads return `[]byte` (not browser `Blob`). Uploads take `fileName`, `size`, and an `io.Reader`, with an optional `onProgress func(float64)` callback.

## Error handling

HTTP status semantics match the API:

- **400** — invalid request / API usage
- **402** — insufficient credits
- **422** — validation errors
- **500** — internal errors

Structured API errors decode into `*maradocs.APIException` (use `errors.As`):

```go
composed, err := client.Pdf.Compose(ctx, req, nil)
if err != nil {
	var apiErr *maradocs.APIException
	if errors.As(err, &apiErr) {
		log.Println(apiErr.Details.Code)    // e.g. 300 (PDF_PAGE_OUT_OF_RANGE)
		log.Println(apiErr.Details.Message)
	}
}
```

Validation endpoints return discriminated responses. Helpers `OkPdf`, `OkImg`, `OkHtml`, `OkEmail`, `OkVideo`, and `OkAudio` return a handle or `*ValidationError` / `*ValidationVirus`.

## API reference

### Server — workspace (`srv.Workspace`)

| Method | Description |
|--------|-------------|
| `Create` | Create a new workspace |
| `Delete` | Delete a workspace |

### Server — account (`srv.Account`)

| Method | Description |
|--------|-------------|
| `GetCredits` | Account-level credit balance |
| `CreateSubaccount` | Create a subaccount |
| `DeleteSubaccount` | Delete a subaccount |
| `GetSubaccountCredits` | Credits for one subaccount |
| `GetSubaccountTransfers` | Reservation or release history |
| `ReserveCredits` / `ReleaseCredits` | Move credits to/from a subaccount |
| `GetSubaccounts` | List subaccounts (optional `startAt`) |

### Server — webview (`srv.Webview`)

| Method | Description |
|--------|-------------|
| `Open` | Open interactive webview session |
| `AddFile` | Add a file to the session |
| `Status` | Session open/closed status |
| `Results` | Files produced in the session |

### PDF (`client.Pdf`)

| Method | Description |
|--------|-------------|
| `Validate` | Validate PDF (virus scan + encoding) |
| `Compose` | Merge/split PDFs by selecting pages |
| `Optimize` | Reduce file size |
| `Rotate` | Rotate specific pages |
| `ToImg` | Render pages as images |
| `Orientation` | Detect and fix page orientation |
| `OcrToPdf` | Searchable PDF with text layer |

### Image (`client.Img`)

| Method | Description |
|--------|-------------|
| `Validate` | Validate image |
| `Thumbnail` | Create thumbnail |
| `FindDocuments` | Detect documents in a photo |
| `ExtractQuadrilateral` | Extract and correct perspective |
| `Orientation` | Detect and fix orientation |
| `Rotate` | Rotate by 0°/90°/180°/270° |
| `ImproveContrast` | Improve contrast |
| `ToJpeg` / `ToPng` / `ToPdf` | Convert format |
| `OcrToPdf` | OCR to searchable PDF |

### HTML (`client.Html`)

| Method | Description |
|--------|-------------|
| `Validate` | Validate HTML |
| `ToPdf` | Convert to PDF |

### Email (`client.Email`)

| Method | Description |
|--------|-------------|
| `Validate` | Parse and validate `.eml` / `.msg`, extract attachments |
| `ToHtml` | Render validated email to HTML |
| `ToPdf` | Render validated email to PDF |

### Video (`client.Video`)

| Method | Description |
|--------|-------------|
| `Validate` | Validate video |

Use `OkVideo` to obtain a `VideoHandle` on success.

### Audio (`client.Audio`)

| Method | Description |
|--------|-------------|
| `Validate` | Validate audio |

Use `OkAudio` to obtain an `AudioHandle` on success.

### Data (`client.Data`)

| Method | Description |
|--------|-------------|
| `CreateUpload` | Mint a **proxy-only** upload capability (`ProxyURL` + `UnvalidatedFileHandle`) |
| `Upload` | Upload file via first-party presigned POST (optional progress callback) |
| `MimeType` | Detect MIME type (async + poll) |
| `VirusScan` | Virus scan (async + poll) |
| `CreateDownloadPdf` / `Jpeg` / `Png` / `Odt` / `Unvalidated` | Mint a **proxy-only** download capability (`ProxyURL`) |
| `DownloadPdf` / `DownloadJpeg` / `DownloadPng` / `DownloadOdt` | Download by handle via first-party SSE-C (optional progress) |
| `DownloadMp4` / `DownloadMp3` / `DownloadWav` / `DownloadFlac` | Transcode and download media (async + poll) |
| `DownloadUnvalidated` | Download unvalidated file (e.g. email body) |

`Create*` methods are for unauthenticated third parties — they always mint `ProxyURL` and never expose SSE-C fields. `Upload` / `Download*` are first-party only:

```go
created, err := client.Data.CreateUpload(ctx, maradocs.DataUploadRequest{
	Size: size,
	Name: &fileName,
})
// Third party: PUT created.ProxyURL with raw body and Content-Length == size
// Integrator: validate(created.UnvalidatedFileHandle) after upload

dl, err := client.Data.CreateDownloadPdf(ctx, maradocs.DataDownloadPdfRequest{
	PdfHandle: pdfHandle,
})
// Third party: GET dl.ProxyURL (no SSE-C headers)
```

### Flow (`client.Flow`)

| Method | Description |
|--------|-------------|
| `OcrImg` / `OcrImgHandle` | Image → searchable PDF pipeline |
| `OcrPdf` / `OcrPdfHandle` | PDF → searchable PDF pipeline |

## Validation helpers

Uploaded files must be validated before use. Poll the validate endpoint, then use the `Ok*` helpers:

```go
up, err := client.Data.Upload(ctx, "doc.pdf", size, reader, nil)
if err != nil {
	return err
}
validated, err := client.Pdf.Validate(ctx, maradocs.PdfValidateRequest{
	UnvalidatedFileHandle: up.UnvalidatedFileHandle,
}, nil)
if err != nil {
	return err
}
pdfHandle, err := maradocs.OkPdf(*validated)
if err != nil {
	return err // *ValidationError or *ValidationVirus
}
```

Same pattern for `OkImg`, `OkHtml`, `OkEmail`, `OkVideo`, and `OkAudio`.

## Examples

### Merge and split PDFs

```go
up1, _ := client.Data.Upload(ctx, "a.pdf", size1, r1, nil)
up2, _ := client.Data.Upload(ctx, "b.pdf", size2, r2, nil)

v1, _ := client.Pdf.Validate(ctx, maradocs.PdfValidateRequest{
	UnvalidatedFileHandle: up1.UnvalidatedFileHandle,
}, nil)
v2, _ := client.Pdf.Validate(ctx, maradocs.PdfValidateRequest{
	UnvalidatedFileHandle: up2.UnvalidatedFileHandle,
}, nil)
pdf1, _ := maradocs.OkPdf(*v1)
pdf2, _ := maradocs.OkPdf(*v2)

pages := []maradocs.PdfComposePdfPage{{PageNumber: 0}, {PageNumber: 2}}
composed, err := client.Pdf.Compose(ctx, maradocs.PdfComposeRequest{
	Pdfs: []maradocs.PdfComposePdf{
		{PdfHandle: pdf1, Pages: &pages},
		{PdfHandle: pdf2}, // nil Pages = all pages
	},
}, nil)
if err != nil {
	return err
}

merged, err := client.Data.DownloadPdf(ctx, maradocs.DataDownloadPdfRequest{
	PdfHandle: composed.PdfHandle,
}, nil)
```

For `PdfComposePdf.Pages`, use a **pointer** to a slice: `nil` omits the field (all pages); `&[]PdfComposePdfPage{}` sends an explicit empty page list. A plain `[]` with `omitempty` would not serialize an empty list correctly.

### Validate and download video

```go
up, _ := client.Data.Upload(ctx, "clip.mp4", size, reader, nil)
validated, _ := client.Video.Validate(ctx, maradocs.VideoValidateRequest{
	UnvalidatedFileHandle: up.UnvalidatedFileHandle,
}, nil)
videoHandle, err := maradocs.OkVideo(*validated)
if err != nil {
	return err
}

mp4, err := client.Data.DownloadMp4(ctx, maradocs.DataDownloadMp4Request{
	VideoHandle: videoHandle,
}, func(p float64) { log.Printf("download %.0f%%", p) }, nil)
```

### Validate and download audio

```go
validated, _ := client.Audio.Validate(ctx, maradocs.AudioValidateRequest{
	UnvalidatedFileHandle: up.UnvalidatedFileHandle,
}, nil)
audioHandle, _ := maradocs.OkAudio(*validated)

mp3, _ := client.Data.DownloadMp3(ctx, maradocs.DataDownloadMp3Request{
	AudioHandle: audioHandle,
}, nil, nil)
wav, _ := client.Data.DownloadWav(ctx, maradocs.DataDownloadWavRequest{
	AudioHandle: audioHandle,
}, nil, nil)
```

### Low-level image processing

```go
up, _ := client.Data.Upload(ctx, "photo.jpg", size, reader, nil)
validated, _ := client.Img.Validate(ctx, maradocs.ImgValidateRequest{
	UnvalidatedFileHandle: up.UnvalidatedFileHandle,
}, nil)
if validated.Response.ClassName != "ImgValidateResponseOk" {
	return fmt.Errorf("validation failed")
}
imgHandle := *validated.Response.ImgHandle

docs, _ := client.Img.FindDocuments(ctx, maradocs.ImgFindDocumentsRequest{
	ImgHandle: imgHandle,
}, nil)
if len(docs.Documents) > 0 {
	extracted, _ := client.Img.ExtractQuadrilateral(ctx, maradocs.ImgExtractQuadrilateralRequest{
		ImgHandle:     imgHandle,
		Quadrilateral: docs.Documents[0].Quadrilateral,
	}, nil)
	_ = extracted.ImgHandle
}
```
