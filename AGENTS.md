# maradocs-sdk-go

Official Go client for the [MaraDocs API](https://api.maradocs.io). Module `github.com/maramia/maradocs-sdk-go`, package `maradocs`. Requires Go 1.22+. Standard library only.

See [README.md](README.md) for method tables and examples. API reference: [api.maradocs.io](https://api.maradocs.io).

## Install

```bash
go get github.com/maramia/maradocs-sdk-go
```

## Clients and auth

| Constructor | Auth | Use |
|-------------|------|-----|
| `NewMaraDocsServer` | Account secret key | Account, workspace, and webview management. Server-side only. |
| `NewMaraDocsClient` | Workspace secret | Document processing (`Data`, `Img`, `Pdf`, `Html`, `Email`, `Video`, `Audio`, `Flow`, `Healthcheck`). |
| `NewHealthcheckClient` | None | Unauthenticated `GET /healthcheck/ping`. |

Never embed the account secret in client-side code. Create a workspace with `MaraDocsServer` and pass `WorkspaceSecret` to `NewMaraDocsClient`.

```go
srv := maradocs.NewMaraDocsServer(maradocs.MaraDocsServerOptions{
	SecretKey: os.Getenv("MARADOCS_SECRET_KEY"),
})
ws, err := srv.Workspace.Create(ctx, maradocs.WorkspaceCreateRequest{})
client, err := maradocs.NewMaraDocsClient(maradocs.MaraDocsClientOptions{
	WorkspaceSecret: ws.WorkspaceSecret,
})
```

Every API method takes `context.Context` as the first argument. Optional `TimeoutMs` on client options is the default timeout (milliseconds) for all requests. Override the API host with `APIURLWithVersion` (default `https://api.maradocs.io/v1`). After `NewMaraDocsClient`, `client.Info` holds decoded workspace metadata.

## Core workflow

Upload → validate → extract a handle → transform → download. `Flow.OcrImg` / `Flow.OcrPdf` skip the low-level steps.

```go
up, err := client.Data.Upload(ctx, "doc.pdf", size, reader, nil)
validated, err := client.Pdf.Validate(ctx, maradocs.PdfValidateRequest{
	UnvalidatedFileHandle: up.UnvalidatedFileHandle,
}, nil)
pdfHandle, err := maradocs.OkPdf(*validated) // *ValidationError or *ValidationVirus
pdfBytes, err := client.Data.DownloadPdf(ctx, maradocs.DataDownloadPdfRequest{
	PdfHandle: pdfHandle,
}, nil)
```

Use `OkPdf`, `OkImg`, `OkHtml`, `OkEmail`, `OkVideo`, `OkAudio` the same way.

## Conventions

- Poll-based methods take `*RequestOptions` as the last argument (`Timeout` in milliseconds). `nil` uses the constructor `TimeoutMs`.
- `CreateUpload` / `CreateDownload*` mint **proxy-only** URLs for unauthenticated third parties. `Upload` / `Download*` are first-party (SSE-C).
- Structured API errors decode into `*maradocs.APIException` (use `errors.As`) with `Details.Code` and `Details.Message`. HTTP: 400 bad request, 402 insufficient credits, 422 validation, 500 internal.
- Exported types and methods are PascalCase; JSON tags are snake_case.
- Downloads return `[]byte`. Uploads take `fileName`, `size`, and `io.Reader`, plus an optional `onProgress func(float64)`.
- For optional lists such as `PdfComposePdf.Pages`, use a **pointer** to a slice: `nil` omits the field (all pages); `&[]PdfComposePdfPage{}` sends an explicit empty list.
