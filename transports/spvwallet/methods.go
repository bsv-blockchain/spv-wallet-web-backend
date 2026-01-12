package spvwallet

import (
	"math"

	"github.com/bsv-blockchain/spv-wallet/models/response"
)

// GetPaymailsFromMetadata returns sender and receiver paymails from metadata.
// If no paymail was found in metadata, fallback paymail is returned.
func GetPaymailsFromMetadata(transaction *response.Transaction, fallbackPaymail string) (string, string) {
	senderPaymail := ""
	receiverPaymail := ""

	if transaction == nil {
		return senderPaymail, receiverPaymail
	}

	if transaction.Metadata != nil {
		// Try to get paymails from metadata if the transaction was made in SPV Wallet.
		if transaction.Metadata["sender"] != nil {
			senderPaymail = transaction.Metadata["sender"].(string)
		}
		if transaction.Metadata["receiver"] != nil {
			receiverPaymail = transaction.Metadata["receiver"].(string)
		}

		if senderPaymail == "" {
			// Try to get paymails from metadata if the transaction was made outside SPV Wallet.
			if transaction.Metadata["p2p_tx_metadata"] != nil {
				p2pTxMetadata := transaction.Metadata["p2p_tx_metadata"].(map[string]interface{})
				if p2pTxMetadata["sender"] != nil {
					senderPaymail = p2pTxMetadata["sender"].(string)
				}
			}
		}
	}

	if transaction.TransactionDirection == "incoming" && receiverPaymail == "" {
		receiverPaymail = fallbackPaymail
	} else if transaction.TransactionDirection == "outgoing" && senderPaymail == "" {
		senderPaymail = fallbackPaymail
	}

	return senderPaymail, receiverPaymail
}

func getAbsoluteValue(value int64) uint64 {
	return uint64(math.Abs(float64(value)))
}
