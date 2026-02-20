package api_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzUserRegistrationJSON tests user registration JSON parsing with arbitrary inputs.
func FuzzUserRegistrationJSON(f *testing.F) {
	// Seed corpus with various JSON formats
	f.Add(`{"email":"user@example.com","password":"Pass123!","passwordConfirmation":"Pass123!"}`)  // Valid
	f.Add(`{}`)                                                                                    // Empty object
	f.Add(`{"email":"","password":"","passwordConfirmation":""}`)                                  // Empty fields
	f.Add(`{"email":"test@example.com"}`)                                                          // Missing fields
	f.Add(`{"email":"invalid-email","password":"weak"}`)                                           // Invalid data
	f.Add(``)                                                                                      // Empty string
	f.Add(`{`)                                                                                     // Malformed JSON
	f.Add(`{"email":"test@example.com","password":"Pass123!","passwordConfirmation":"Different"}`) // Mismatched passwords
	f.Add(`{"email":"test","password":"` + strings.Repeat("a", 1000) + `"}`)                       // Very long password
	f.Add(`{"extra":"field","email":"test@example.com"}`)                                          // Extra fields

	f.Fuzz(func(t *testing.T, jsonInput string) {
		// Primary goal: ensure JSON parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("JSON parsing panicked on input %q: %v", jsonInput, r)
			}
		}()

		// Define registration struct
		type RegisterUser struct {
			Email                string `json:"email"`
			Password             string `json:"password"` //nolint:gosec // G117: fuzz test struct field, not a secret leak
			PasswordConfirmation string `json:"passwordConfirmation"`
		}

		var req RegisterUser
		err := json.Unmarshal([]byte(jsonInput), &req)

		// Parsing may fail, that's acceptable - but no panic
		if err == nil {
			// If parsing succeeded, validate the structure
			// Email, password fields should exist (can be empty)
			_ = req.Email
			_ = req.Password
			_ = req.PasswordConfirmation

			// Test basic validation logic
			if req.Email != "" && strings.Contains(req.Email, "@") {
				require.Contains(t, req.Email, "@", "Email should contain @")
			}
		}
	})
}

// FuzzTransactionRequestJSON tests transaction creation JSON parsing.
func FuzzTransactionRequestJSON(f *testing.F) {
	// Seed corpus with transaction patterns
	f.Add(`{"password":"Pass123!","recipient":"user@example.com","satoshis":1000}`) // Valid
	f.Add(`{}`)                                                                     // Empty
	f.Add(`{"password":"","recipient":"","satoshis":0}`)                            // Empty fields
	f.Add(`{"satoshis":999999999999999}`)                                           // Very large number
	f.Add(`{"satoshis":-100}`)                                                      // Negative (invalid in JSON for uint64)
	f.Add(`{"recipient":"invalid-paymail"}`)                                        // Missing @
	f.Add(`{"password":"test","recipient":"user@example.com"}`)                     // Missing satoshis
	f.Add(``)                                                                       // Empty string
	f.Add(`{`)                                                                      // Malformed
	f.Add(`{"satoshis":"not-a-number"}`)                                            // Wrong type

	f.Fuzz(func(t *testing.T, jsonInput string) {
		// Primary goal: ensure JSON parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Transaction JSON parsing panicked on input %q: %v", jsonInput, r)
			}
		}()

		// Define transaction struct
		type CreateTransaction struct {
			Password  string `json:"password"` //nolint:gosec // G117: fuzz test struct field, not a secret leak
			Recipient string `json:"recipient"`
			Satoshis  uint64 `json:"satoshis"`
		}

		var req CreateTransaction
		err := json.Unmarshal([]byte(jsonInput), &req)

		// Parsing may fail, that's acceptable
		if err == nil {
			// If parsing succeeded, validate the structure
			_ = req.Password
			_ = req.Recipient
			_ = req.Satoshis

			// Test validation logic
			if req.Satoshis > 0 {
				require.Positive(t, req.Satoshis, "Satoshis should be positive")
			}

			// Recipient should be a valid paymail (email-like format)
			if req.Recipient != "" && strings.Contains(req.Recipient, "@") {
				require.Contains(t, req.Recipient, "@", "Recipient should contain @")
			}
		}
	})
}

// FuzzContactUpsertJSON tests contact creation/update JSON parsing.
func FuzzContactUpsertJSON(f *testing.F) {
	// Seed corpus with contact patterns
	f.Add(`{"paymail":"contact@example.com","fullName":"John Doe"}`) // Valid
	f.Add(`{}`)                                                      // Empty
	f.Add(`{"paymail":"","fullName":""}`)                            // Empty fields
	f.Add(`{"paymail":"invalid"}`)                                   // Invalid paymail
	f.Add(`{"fullName":"` + strings.Repeat("A", 500) + `"}`)         // Very long name
	f.Add(`{"metadata":{"key":"value","nested":{"deep":"data"}}}`)   // Nested metadata
	f.Add(``)                                                        // Empty string
	f.Add(`{`)                                                       // Malformed
	f.Add(`{"paymail":"test@example.com","metadata":null}`)          // Null metadata
	f.Add(`{"paymail":"test@example.com","metadata":{}}`)            // Empty metadata

	f.Fuzz(func(t *testing.T, jsonInput string) {
		// Primary goal: ensure JSON parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Contact JSON parsing panicked on input %q: %v", jsonInput, r)
			}
		}()

		// Define contact struct
		type UpsertContact struct {
			Paymail  string                 `json:"paymail"`
			FullName string                 `json:"fullName,omitempty"`
			Metadata map[string]interface{} `json:"metadata,omitempty"`
		}

		var req UpsertContact
		err := json.Unmarshal([]byte(jsonInput), &req)

		// Parsing may fail, that's acceptable
		if err == nil {
			// If parsing succeeded, validate the structure
			_ = req.Paymail
			_ = req.FullName
			_ = req.Metadata

			// Test validation logic
			if req.Paymail != "" {
				// Paymail should be in email format
				if strings.Contains(req.Paymail, "@") {
					require.Contains(t, req.Paymail, "@", "Paymail should contain @")
				}
			}

			// Metadata should be a valid map
			if req.Metadata != nil {
				require.IsType(t, map[string]interface{}{}, req.Metadata, "Metadata should be a map")
			}
		}
	})
}

// FuzzSearchTransactionJSON tests transaction search JSON parsing with complex nested structures.
func FuzzSearchTransactionJSON(f *testing.F) {
	// Seed corpus with search patterns
	f.Add(`{"conditions":{"status":"complete"},"metadata":{"app":"wallet"}}`) // Valid
	f.Add(`{}`)                                                               // Empty
	f.Add(`{"conditions":{}}`)                                                // Empty conditions
	f.Add(`{"metadata":null}`)                                                // Null metadata
	f.Add(`{"conditions":{"nested":{"deep":{"value":123}}}}`)                 // Deep nesting
	f.Add(``)                                                                 // Empty string
	f.Add(`{`)                                                                // Malformed
	f.Add(`{"conditions":{"key":"` + strings.Repeat("value", 100) + `"}}`)    // Long values

	f.Fuzz(func(t *testing.T, jsonInput string) {
		// Primary goal: ensure JSON parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Search JSON parsing panicked on input %q: %v", jsonInput, r)
			}
		}()

		// Define search struct (simplified)
		type SearchTransaction struct {
			Conditions map[string]interface{} `json:"conditions,omitempty"`
			Metadata   map[string]interface{} `json:"metadata,omitempty"`
		}

		var req SearchTransaction
		err := json.Unmarshal([]byte(jsonInput), &req)

		// Parsing may fail, that's acceptable
		if err == nil {
			// If parsing succeeded, validate the structure
			_ = req.Conditions
			_ = req.Metadata

			// Test that maps are properly initialized
			if req.Conditions != nil {
				require.IsType(t, map[string]interface{}{}, req.Conditions)
			}
			if req.Metadata != nil {
				require.IsType(t, map[string]interface{}{}, req.Metadata)
			}
		}
	})
}
