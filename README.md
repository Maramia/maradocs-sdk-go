# maradocs-sdk-go

Go client for the [MaraDocs](https://maradocs.io) API. This module mirrors the public TypeScript SDK (`maradocs-sdk-ts`): `MaraDocsClient` (workspace secret), `MaraDocsServer` (account secret), polling behaviour, and the high-level `Flow` helper.

## Install

```bash
go get github.com/maramia/maradocs-sdk-go
```

## Quick start

```go
ctx := context.Background()
client, err := maradocs.NewMaraDocsClient(maradocs.MaraDocsClientOptions{
    WorkspaceSecret:   os.Getenv("MARADOCS_WORKSPACE_SECRET"),
    APIURLWithVersion: "", // default https://api.maradocs.io/v1
})
if err != nil {
    log.Fatal(err)
}
// client.Info holds decoded workspace metadata (account id, workspace id, encryption key bytes).
_, _ = client.Healthcheck.Ping(ctx)
```

Server-side workspace lifecycle:

```go
srv := maradocs.NewMaraDocsServer(maradocs.MaraDocsServerOptions{
    SecretKey: os.Getenv("MARADOCS_SECRET_KEY"),
})
ws, err := srv.Workspace.Create(ctx, maradocs.WorkspaceCreateRequest{})
```

## API surface

- **Client**: `Healthcheck`, `Data` (upload, MIME, virus scan, downloads), `Img`, `Pdf`, `Video`, `Audio`, `Html`, `Email`, `Flow`.
- **Server**: `Healthcheck`, `Account`, `Workspace`, `Webview`.
- **Options**: poll-based methods accept `*maradocs.RequestOptions{ Timeout: &ms }` to override per-call timeouts (mirrors TS `RequestOptions`).

Binary uploads return raw `[]byte` from download helpers instead of `Blob`.
