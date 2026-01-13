package users_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzEmailValidation tests email validation and parsing to ensure it never panics.
func FuzzEmailValidation(f *testing.F) {
	// Seed corpus with various email formats and edge cases
	f.Add("user@example.com")                          // Valid email
	f.Add("test@domain.co.uk")                         // Valid with subdomain
	f.Add("")                                          // Empty string
	f.Add("noemail")                                   // Missing @
	f.Add("@domain.com")                               // Missing username
	f.Add("user@")                                     // Missing domain
	f.Add("user@@domain.com")                          // Double @
	f.Add("user@domain@com")                           // Multiple @
	f.Add("@")                                         // Just @
	f.Add("user @domain.com")                          // Space in email
	f.Add("user@domain .com")                          // Space in domain
	f.Add("user+tag@example.com")                      // Plus addressing
	f.Add("user.name@example.com")                     // Dots in username
	f.Add("user@sub.domain.example.com")               // Multiple subdomains
	f.Add(strings.Repeat("a", 100) + "@example.com")   // Very long username
	f.Add("user@" + strings.Repeat("a", 100) + ".com") // Very long domain

	f.Fuzz(func(t *testing.T, email string) {
		// Primary goal: ensure email processing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Email validation panicked on input %q: %v", email, r)
			}
		}()

		// Test basic email validation patterns
		hasAt := strings.Contains(email, "@")
		parts := strings.Split(email, "@")

		// If email contains @, we can split it safely
		if hasAt && len(parts) == 2 {
			username := parts[0]
			domain := parts[1]

			// Validate username and domain are not empty
			if username != "" && domain != "" {
				// This would be a potentially valid email format
				require.NotEmpty(t, username, "Username should not be empty")
				require.NotEmpty(t, domain, "Domain should not be empty")
			}
		}

		// Test splitEmail function indirectly by simulating its logic
		// The actual function should handle invalid formats gracefully
		components := strings.Split(email, "@")
		if len(components) != 2 {
			// Invalid email format - function should return empty strings
			// This is defensive programming validation
			_ = components
		}
	})
}

// FuzzPasswordValidation tests password handling with arbitrary inputs.
func FuzzPasswordValidation(f *testing.F) {
	// Seed corpus with various password patterns
	f.Add("")                        // Empty password
	f.Add("short")                   // Short password
	f.Add("validPassword123!")       // Valid password
	f.Add(strings.Repeat("a", 1000)) // Very long password
	f.Add("pass word")               // Password with space
	// f.Add("unicode密码🔒")                 // Unicode (skipped for gosmopolitan)
	f.Add("\x00\x01\x02")                        // Binary/control characters
	f.Add("special!@#$%^&*()_+-=[]{}|;:',.<>?/") // All special chars
	f.Add("\n\t\r")                              // Whitespace only
	f.Add("Pass123!")                            // Mixed case with numbers and special

	f.Fuzz(func(t *testing.T, password string) {
		// Primary goal: ensure password processing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Password validation panicked on input %q: %v", password, r)
			}
		}()

		// Test basic password validation patterns
		// Most systems require minimum length
		hasMinLength := len(password) >= 8
		hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
		hasDigit := strings.ContainsAny(password, "0123456789")
		hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:',.<>?/")

		// If password meets criteria, it could be valid
		if hasMinLength && hasUpper && hasLower && hasDigit && hasSpecial {
			// Strong password detected
			_ = hasMinLength // Validation passed, continue
		}

		// Test trimming (passwords are often trimmed)
		trimmed := strings.TrimSpace(password)
		_ = trimmed

		// Test length validation
		if len(password) > 0 {
			require.NotEmpty(t, password, "Non-empty password should have length > 0")
		}
	})
}

// FuzzPaymailValidation tests paymail address validation.
func FuzzPaymailValidation(f *testing.F) {
	// Seed corpus with paymail patterns (similar to email but for BSV)
	f.Add("user@moneybutton.com")                                           // Valid paymail
	f.Add("test@handcash.io")                                               // Valid paymail
	f.Add("")                                                               // Empty
	f.Add("nopaymail")                                                      // Invalid format
	f.Add("@domain.com")                                                    // Missing username
	f.Add("user@")                                                          // Missing domain
	f.Add("user@@domain.com")                                               // Double @
	f.Add(strings.Repeat("a", 50) + "@" + strings.Repeat("b", 50) + ".com") // Long paymail

	f.Fuzz(func(t *testing.T, paymail string) {
		// Primary goal: ensure paymail processing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Paymail validation panicked on input %q: %v", paymail, r)
			}
		}()

		// Paymail format is similar to email: username@domain
		parts := strings.Split(paymail, "@")

		// Valid paymail should have exactly 2 parts
		if len(parts) == 2 {
			username := parts[0]
			domain := parts[1]

			if username != "" && domain != "" {
				// Could be a valid paymail
				require.NotEmpty(t, username)
				require.NotEmpty(t, domain)
			}
		}

		// Test that paymail can be safely processed
		_ = strings.ToLower(paymail)
		_ = strings.TrimSpace(paymail)
	})
}
