package transactions

import (
	"github.com/bsv-blockchain/spv-wallet/models"
	"github.com/bsv-blockchain/spv-wallet/models/filter"
)

// CreateTransaction represents request for creating new transaction.
type CreateTransaction struct {
	Password  string `json:"password"` //nolint:gosec // G117: field name matches pattern but this is a transaction request model, not a secret leak
	Recipient string `json:"recipient"`
	Satoshis  uint64 `json:"satoshis"`
}

// SearchTransaction represents request for searching transactions.
type SearchTransaction struct {
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    models.Metadata        `json:"metadata,omitempty"`
	QueryParams *filter.QueryParams    `json:"params,omitempty"`
}
