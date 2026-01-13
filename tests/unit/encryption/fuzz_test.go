package encryption_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bsv-blockchain/spv-wallet-web-backend/encryption"
)

// FuzzDecrypt tests the Decrypt function with arbitrary ciphertext inputs to ensure it never panics.
func FuzzDecrypt(f *testing.F) {
	// Seed corpus with edge cases
	f.Add("")                                    // Empty string
	f.Add("no-dashes")                           // Missing delimiters
	f.Add("one-dash")                            // Only one dash
	f.Add("too-many-dashes-here-and-there-more") // Too many dashes
	f.Add("invalid-hex-chars")                   // Invalid hex characters
	f.Add("abc123-def456-ghi789")                // Invalid hex data
	f.Add("--")                                  // Only dashes
	f.Add("a-b-c")                               // Too short components

	// Add a valid encrypted string from the test suite
	validEncrypted, err := encryption.Encrypt("test-passphrase", "test data")
	require.NoError(f, err)
	f.Add(validEncrypted)

	f.Fuzz(func(t *testing.T, ciphertext string) {
		// Primary goal: function should NEVER panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Decrypt panicked on input %q: %v", ciphertext, r)
			}
		}()

		// Call Decrypt - it should return empty string on invalid input, never panic
		result := encryption.Decrypt("test-passphrase", ciphertext)

		// Result can be empty or garbage for invalid input, that's acceptable
		// We just verify no panic occurred
		_ = result
	})
}

// FuzzEncryptDecryptRoundTrip tests encryption and decryption round-trip with arbitrary inputs.
func FuzzEncryptDecryptRoundTrip(f *testing.F) {
	// Seed corpus with various edge cases
	f.Add("passphrase123", "")                          // Empty plaintext
	f.Add("test", "simple text")                        // Simple text
	f.Add("key", "text with special chars: !@#$%^&*()") // Special characters
	// f.Add("pass", "unicode: 你好世界 🚀")                 // Unicode (skipped for gosmopolitan)
	f.Add("secret", strings.Repeat("a", 1000))     // Long text (1KB)
	f.Add("pwd", "\n\t\r special whitespace")      // Whitespace characters
	f.Add("test123", "Lorem ipsum dolor sit amet") // Normal sentence
	f.Add("", "empty passphrase test")             // Empty passphrase
	f.Add("p@ss!", "text")                         // Special chars in passphrase

	f.Fuzz(func(t *testing.T, passphrase, plaintext string) {
		// Primary goal: no panics during encryption or decryption
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Round-trip panicked with passphrase=%q plaintext=%q: %v", passphrase, plaintext, r)
			}
		}()

		// Encrypt the plaintext
		encrypted, err := encryption.Encrypt(passphrase, plaintext)
		if err != nil {
			// Encryption can fail on certain inputs, that's acceptable
			return
		}

		// Decrypt should recover the original plaintext with correct passphrase
		decrypted := encryption.Decrypt(passphrase, encrypted)
		require.Equal(t, plaintext, decrypted, "Round-trip failed: expected %q, got %q", plaintext, decrypted)

		// Decrypt with wrong passphrase should NOT recover original (security check)
		wrongDecrypted := encryption.Decrypt("wrong-passphrase-12345", encrypted)
		if plaintext != "" {
			require.NotEqual(t, plaintext, wrongDecrypted, "Decryption with wrong passphrase should not recover plaintext")
		}
	})
}

// FuzzHash tests the Hash function with arbitrary inputs to ensure consistent behavior.
func FuzzHash(f *testing.F) {
	// Seed corpus with edge cases
	f.Add("")                  // Empty string
	f.Add("test")              // Simple string
	f.Add("special!@#$%^&*()") // Special characters
	// f.Add("unicode: 你好世界 🚀")      // Unicode (skipped for gosmopolitan)
	f.Add(strings.Repeat("a", 10000))   // Very long string (10KB)
	f.Add("\x00\x01\x02\x03")           // Binary data
	f.Add("multi\nline\ntext")          // Multiline
	f.Add("Lorem ipsum dolor sit amet") // Normal text

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: no panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Hash panicked on input %q: %v", input, r)
			}
		}()

		// Hash the input
		hash, err := encryption.Hash(input)
		require.NoError(t, err, "Hash should not return error")

		// Validate hash properties
		require.NotEmpty(t, hash, "Hash should not be empty")
		require.Len(t, hash, 64, "SHA256 hash should be 64 hex characters")

		// Verify hash is deterministic (same input = same output)
		hash2, err2 := encryption.Hash(input)
		require.NoError(t, err2)
		require.Equal(t, hash, hash2, "Hash should be deterministic")

		// Verify hash only contains valid hex characters
		for _, c := range hash {
			require.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'), "Hash should only contain hex characters, got %q", hash)
		}
	})
}
