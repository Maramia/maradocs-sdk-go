package maradocs

// Subaccount is returned by account endpoints.
type Subaccount struct {
	ID               string `json:"id"`
	AvailableCredits int    `json:"available_credits"`
	CreatedAt        string `json:"created_at"`
}

// AccountCreditsGetResponse lists account-level credits.
type AccountCreditsGetResponse struct {
	Credits int `json:"credits"`
}

// SubaccountOverview lists subaccounts.
type SubaccountOverview struct {
	Subaccounts []Subaccount `json:"subaccounts"`
}

// SubaccountTransferType is reservation or release.
type SubaccountTransferType string

const (
	SubaccountTransferReservation SubaccountTransferType = "reservation"
	SubaccountTransferRelease     SubaccountTransferType = "release"
)

// SubaccountTransfer is one credit transfer record.
type SubaccountTransfer struct {
	Amount       int                    `json:"amount"`
	CreatedAt    string                 `json:"created_at"`
	TransferType SubaccountTransferType `json:"transfer_type"`
}

// SubaccountTransferHistory lists transfers for a subaccount.
type SubaccountTransferHistory struct {
	Transfers []SubaccountTransfer `json:"transfers"`
}

// AccountSecret is decoded account secret material (not used by all endpoints).
type AccountSecret struct {
	AccountID string `json:"account_id"`
	Seed      string `json:"seed"`
}
