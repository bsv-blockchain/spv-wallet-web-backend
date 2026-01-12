package data

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/bsv-blockchain/spv-wallet-web-backend/transports/spvwallet"
)

// CreateTestTransactions returns 'count' randomly generated transactions.
func CreateTestTransactions(count int) []spvwallet.FullTransaction {
	result := make([]spvwallet.FullTransaction, count)
	gofakeit.Slice(&result)

	return result
}
