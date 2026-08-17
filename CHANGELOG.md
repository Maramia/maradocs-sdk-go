# Changelog

## 0.2.0

- Optional `use_proxy` / `proxy_url` on data upload and download models for the SSE-C encrypting/decrypting proxy.
- `DataEp.CreateUpload` and `CreateDownloadPdf` / `Jpeg` / `Png` / `Odt` / `Unvalidated` mint **proxy-only** capability URLs (`CreateUploadResult` / `CreateDownloadProxyResult`). They never return SSE-C `post_url`/`post_header` or `url`/`headers`. Share only `ProxyURL` with unauthenticated third parties.
- First-party `Upload` / `Download*` strip `UseProxy` if present.

## 0.1.0

- Initial release mirroring the public TypeScript SDK surface.
