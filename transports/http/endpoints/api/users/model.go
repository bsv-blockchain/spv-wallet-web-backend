package users

import "github.com/bsv-blockchain/spv-wallet-web-backend/domain/users"

// RegisterUser is a struct that contains user register data.
type RegisterUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"` //nolint:gosec // G117: field name matches pattern but this is a registration request model, not a secret leak
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// RegisterResponse represents response that is sent after user creation.
type RegisterResponse struct {
	Mnemonic string `json:"mnemonic"`
	Paymail  string `json:"paymail"`
}

// UserResponse is a struct that represents user information.
type UserResponse struct {
	UserID  int           `json:"userId"`
	Paymail string        `json:"paymail"`
	Email   string        `json:"email"`
	Balance users.Balance `json:"balance"`
}
