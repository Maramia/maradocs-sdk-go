package maradocs

import (
	"fmt"
	"strings"
)

// MaraDocsClientOptions configures a workspace-scoped MaraDocs client.
type MaraDocsClientOptions struct {
	WorkspaceSecret   string
	APIURLWithVersion string
	TimeoutMs         *int
}

// MaraDocsClient is the public workspace client (mirrors TS MaraDocsClient).
type MaraDocsClient struct {
	Info WorkspaceInfo

	Healthcheck *HealthcheckEp
	Data        *DataEp
	Img         *ImgEp
	Pdf         *PdfEp
	Video       *VideoEp
	Audio       *AudioEp
	Html        *HtmlEp
	Email       *EmailEp
	Flow        *Flow
}

// NewMaraDocsClient constructs a client for a workspace secret.
func NewMaraDocsClient(opt MaraDocsClientOptions) (*MaraDocsClient, error) {
	info, err := DecodeWorkspaceInfo(opt.WorkspaceSecret)
	if err != nil {
		return nil, fmt.Errorf("workspace secret: %w", err)
	}
	api := strings.TrimSpace(opt.APIURLWithVersion)
	if api == "" {
		api = defaultAPIURLV1
	}
	t := newTransport(opt.WorkspaceSecret, api, opt.TimeoutMs)
	data := newDataEp(t)
	img := newImgEp(t)
	pdf := newPdfEp(t)
	return &MaraDocsClient{
		Info:        *info,
		Healthcheck: newHealthcheckEp(api, opt.TimeoutMs),
		Data:        data,
		Img:         img,
		Pdf:         pdf,
		Video:       newVideoEp(t),
		Audio:       newAudioEp(t),
		Html:        newHtmlEp(t),
		Email:       newEmailEp(t),
		Flow:        newFlow(data, img, pdf),
	}, nil
}

// MaraDocsServerOptions configures an account-secret server client.
type MaraDocsServerOptions struct {
	SecretKey         string
	APIURLWithVersion string
	TimeoutMs         *int
}

// MaraDocsServer exposes account/workspace/webview operations (mirrors TS MaraDocsServer).
type MaraDocsServer struct {
	Healthcheck *HealthcheckEp
	Account     *AccountEp
	Workspace   *WorkspaceEp
	Webview     *WebviewEp

	t *transport
}

// NewMaraDocsServer constructs a server-side client using the account secret key.
func NewMaraDocsServer(opt MaraDocsServerOptions) *MaraDocsServer {
	api := strings.TrimSpace(opt.APIURLWithVersion)
	if api == "" {
		api = defaultAPIURLV1
	}
	t := newTransport(opt.SecretKey, api, opt.TimeoutMs)
	return &MaraDocsServer{
		Healthcheck: newHealthcheckEp(api, opt.TimeoutMs),
		Account:     newAccountEp(t),
		Workspace:   newWorkspaceEp(t),
		Webview:     newWebviewEp(t),
		t:           t,
	}
}

// NewHealthcheckClient constructs an unauthenticated healthcheck helper (GET /healthcheck/ping).
func NewHealthcheckClient(apiURLWithVersion string, timeoutMs *int) *HealthcheckEp {
	api := strings.TrimSpace(apiURLWithVersion)
	if api == "" {
		api = defaultAPIURLV1
	}
	return newHealthcheckEp(api, timeoutMs)
}
