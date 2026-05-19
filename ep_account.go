package maradocs

import (
	"context"
	"fmt"
	"net/url"
)

// AccountEp implements /account/* server-side endpoints.
type AccountEp struct {
	t *transport
}

func newAccountEp(t *transport) *AccountEp {
	return &AccountEp{t: t}
}

// GetCredits returns account-level credits.
func (a *AccountEp) GetCredits(ctx context.Context) (*AccountCreditsGetResponse, error) {
	var out AccountCreditsGetResponse
	if err := a.t.getJSON(ctx, "/account/credits", &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateSubaccount creates an empty subaccount (POST body {}).
func (a *AccountEp) CreateSubaccount(ctx context.Context) (*Subaccount, error) {
	var out Subaccount
	if err := a.t.postJSON(ctx, "/account/subaccount", struct{}{}, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteSubaccount deletes a subaccount by id.
func (a *AccountEp) DeleteSubaccount(ctx context.Context, id string) error {
	path := fmt.Sprintf("/account/subaccount/%s", url.PathEscape(id))
	return a.t.deleteJSON(ctx, path, nil, nil, nil)
}

// GetSubaccountCredits returns credit balance for a subaccount.
func (a *AccountEp) GetSubaccountCredits(ctx context.Context, id string) (*Subaccount, error) {
	var out Subaccount
	path := fmt.Sprintf("/account/subaccount/%s/credits", url.PathEscape(id))
	if err := a.t.getJSON(ctx, path, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// SubaccountTransfersParams configures transfer history listing.
type SubaccountTransfersParams struct {
	ID            string
	TransferType  SubaccountTransferType
	CreatedBefore *string
	Limit         *int
}

// GetSubaccountTransfers returns transfer history for a subaccount.
func (a *AccountEp) GetSubaccountTransfers(ctx context.Context, p SubaccountTransfersParams) (*SubaccountTransferHistory, error) {
	base := fmt.Sprintf("/account/subaccount/%s/%s", url.PathEscape(p.ID), p.TransferType)
	q := url.Values{}
	if p.CreatedBefore != nil {
		q.Set("created_before", *p.CreatedBefore)
	}
	if p.Limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *p.Limit))
	}
	path := base
	if enc := q.Encode(); enc != "" {
		path = base + "?" + enc
	}
	var out SubaccountTransferHistory
	if err := a.t.getJSON(ctx, path, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReserveCredits reserves credits on a subaccount.
func (a *AccountEp) ReserveCredits(ctx context.Context, id string, amount int) error {
	path := fmt.Sprintf("/account/subaccount/%s/reservation/%d", url.PathEscape(id), amount)
	return a.t.patchJSON(ctx, path, nil, nil, nil)
}

// ReleaseCredits releases credits from a subaccount.
func (a *AccountEp) ReleaseCredits(ctx context.Context, id string, amount int) error {
	path := fmt.Sprintf("/account/subaccount/%s/release/%d", url.PathEscape(id), amount)
	return a.t.patchJSON(ctx, path, nil, nil, nil)
}

// GetSubaccounts lists subaccounts (optional pagination via start_at ISO timestamp).
func (a *AccountEp) GetSubaccounts(ctx context.Context, startAt *string) (*SubaccountOverview, error) {
	path := "/account/subaccounts/credits"
	if startAt != nil && *startAt != "" {
		path = path + "?start_at=" + url.QueryEscape(*startAt)
	}
	var out SubaccountOverview
	if err := a.t.getJSON(ctx, path, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
