package access

import "github.com/bsv-blockchain/spv-wallet-web-backend/domain/users"

// SignInUser is a struct that contains user sign in data.
type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"` //nolint:gosec // G117: field name matches pattern but this is an auth request model, not a secret leak
}

// SignInResponse is a struct that represents struct sended after user sign in.
type SignInResponse struct {
	Paymail string        `json:"paymail"`
	Balance users.Balance `json:"balance"`
}
