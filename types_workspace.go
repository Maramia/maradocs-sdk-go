package maradocs

// WorkspaceCreateRequest creates a workspace.
type WorkspaceCreateRequest struct {
	Subaccount *string `json:"subaccount,omitempty"`
	FileTTL    *int    `json:"file_ttl,omitempty"`
}

// WorkspaceCreateResponse contains workspace credentials.
type WorkspaceCreateResponse struct {
	Subaccount      *string `json:"subaccount"`
	WorkspaceID     string  `json:"workspace_id"`
	WorkspaceSecret string  `json:"workspace_secret"`
}

// WorkspaceDeleteRequest deletes a workspace.
type WorkspaceDeleteRequest struct {
	WorkspaceID string  `json:"workspace_id"`
	Subaccount  *string `json:"subaccount,omitempty"`
}

// WorkspaceDeleteResponse is an empty JSON object on success.
type WorkspaceDeleteResponse struct{}
