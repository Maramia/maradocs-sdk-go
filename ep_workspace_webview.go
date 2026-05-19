package maradocs

import "context"

// WorkspaceEp implements /workspace endpoints (server secret).
type WorkspaceEp struct {
	t *transport
}

func newWorkspaceEp(t *transport) *WorkspaceEp {
	return &WorkspaceEp{t: t}
}

// Create creates a new workspace.
func (w *WorkspaceEp) Create(ctx context.Context, req WorkspaceCreateRequest) (*WorkspaceCreateResponse, error) {
	var out WorkspaceCreateResponse
	if err := w.t.postJSON(ctx, "/workspace", req, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a workspace.
func (w *WorkspaceEp) Delete(ctx context.Context, req WorkspaceDeleteRequest) error {
	var out WorkspaceDeleteResponse
	return w.t.deleteJSON(ctx, "/workspace", req, &out, nil)
}

// WebviewEp implements /webview endpoints.
type WebviewEp struct {
	t *transport
}

func newWebviewEp(t *transport) *WebviewEp {
	return &WebviewEp{t: t}
}

// Open opens the interactive webview and returns a URL.
func (w *WebviewEp) Open(ctx context.Context, req WebviewOpenRequest) (*WebviewOpenResponse, error) {
	var out WebviewOpenResponse
	if err := w.t.postJSON(ctx, "/webview", req, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFile adds a file to the webview session.
func (w *WebviewEp) AddFile(ctx context.Context, req WebviewAddFileRequest) (*WebviewAddFileResponse, error) {
	var out WebviewAddFileResponse
	if err := w.t.putJSON(ctx, "/webview", req, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// Status returns webview open/closed status.
func (w *WebviewEp) Status(ctx context.Context) (*WebviewStatusResponse, error) {
	var out WebviewStatusResponse
	if err := w.t.getJSON(ctx, "/webview/status", &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// Results returns files produced in the webview session.
func (w *WebviewEp) Results(ctx context.Context) (*WebviewResultsResponse, error) {
	var out WebviewResultsResponse
	if err := w.t.getJSON(ctx, "/webview/results", &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
