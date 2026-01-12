package transactions

import "github.com/bsv-blockchain/spv-wallet-web-backend/domain/users"

// PaginatedTransactions represents transactions with pagination details
// like transactins count and number of pages.
type PaginatedTransactions struct {
	Count        int64               `json:"count"`
	Pages        int                 `json:"pages"`
	Transactions []users.Transaction `json:"transactions"`
}
