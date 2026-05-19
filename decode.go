package maradocs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// workspaceInfoWire is the JSON payload embedded in the workspace secret (after the first 64 bytes).
type workspaceInfoWire struct {
	AccountID      string  `json:"account_id"`
	Subaccount     *string `json:"subaccount"`
	WorkspaceID    string  `json:"workspace_id"`
	EncryptionKey  string  `json:"encryption_key"`
}

// WorkspaceInfo contains fields decoded from a workspace secret (publishable JWT payload).
type WorkspaceInfo struct {
	AccountID      string
	Subaccount     *string
	WorkspaceID    string
	EncryptionKey  []byte
}

// DecodeWorkspaceInfo decodes the workspace secret the same way as the TypeScript SDK.
func DecodeWorkspaceInfo(token string) (*WorkspaceInfo, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("decode workspace secret: %w", err)
	}
	if len(decoded) <= 64 {
		return nil, fmt.Errorf("workspace secret too short")
	}
	payload := decoded[64:]
	var w workspaceInfoWire
	if err := json.Unmarshal(payload, &w); err != nil {
		return nil, fmt.Errorf("parse workspace payload: %w", err)
	}
	enc, err := base64.StdEncoding.DecodeString(w.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decode encryption_key: %w", err)
	}
	return &WorkspaceInfo{
		AccountID:     w.AccountID,
		Subaccount:    w.Subaccount,
		WorkspaceID:   w.WorkspaceID,
		EncryptionKey: enc,
	}, nil
}
