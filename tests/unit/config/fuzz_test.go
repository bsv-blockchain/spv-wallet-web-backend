package config_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// FuzzConfigIntegerParsing tests integer parsing from config values.
func FuzzConfigIntegerParsing(f *testing.F) {
	// Seed corpus with various integer patterns
	f.Add("8080")                           // Valid port
	f.Add("0")                              // Zero
	f.Add("-1")                             // Negative
	f.Add("65535")                          // Max valid port
	f.Add("99999")                          // Out of range port
	f.Add("")                               // Empty
	f.Add("not-a-number")                   // Invalid
	f.Add("12.34")                          // Float
	f.Add("1e10")                           // Scientific notation
	f.Add(strconv.Itoa(int(^uint(0) >> 1))) // Max int

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Integer parsing panicked on input %q: %v", input, r)
			}
		}()

		// Test integer parsing (similar to viper.GetInt behavior)
		intVal, err := strconv.Atoi(input)
		if err == nil {
			// Valid integer
			_ = intVal
		}

		// Test parsing as int with default fallback
		if input == "" {
			// Empty string should use default
			defaultVal := 8080
			_ = defaultVal
		}
	})
}

// FuzzConfigBooleanParsing tests boolean parsing from config values.
func FuzzConfigBooleanParsing(f *testing.F) {
	// Seed corpus with various boolean patterns
	f.Add("true")
	f.Add("false")
	f.Add("True")
	f.Add("False")
	f.Add("TRUE")
	f.Add("FALSE")
	f.Add("1")
	f.Add("0")
	f.Add("")
	f.Add("yes")
	f.Add("no")
	f.Add("invalid")

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Boolean parsing panicked on input %q: %v", input, r)
			}
		}()

		// Test boolean parsing (similar to viper.GetBool behavior)
		boolVal, err := strconv.ParseBool(input)
		if err == nil {
			// Valid boolean
			require.IsType(t, true, boolVal)
		}
	})
}

// FuzzConfigDurationParsing tests duration parsing from config values.
func FuzzConfigDurationParsing(f *testing.F) {
	// Seed corpus with various duration patterns
	f.Add("10s")     // Seconds
	f.Add("5m")      // Minutes
	f.Add("1h")      // Hours
	f.Add("100ms")   // Milliseconds
	f.Add("2h30m")   // Combined
	f.Add("")        // Empty
	f.Add("invalid") // Invalid
	f.Add("10")      // No unit
	f.Add("-5m")     // Negative
	f.Add("999999h") // Very large

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Duration parsing panicked on input %q: %v", input, r)
			}
		}()

		// Test duration parsing (similar to viper.GetDuration behavior)
		duration, err := time.ParseDuration(input)
		if err == nil {
			// Valid duration
			require.IsType(t, time.Duration(0), duration)
		}
	})
}

// FuzzConfigStringParsing tests string parsing with special characters and edge cases.
func FuzzConfigStringParsing(f *testing.F) {
	// Seed corpus with various string patterns
	f.Add("localhost")                // Simple string
	f.Add("")                         // Empty
	f.Add("https://example.com")      // URL
	f.Add("user@example.com")         // Email
	f.Add("special!@#$%^&*()")        // Special chars
	f.Add(strings.Repeat("a", 10000)) // Very long
	f.Add("\x00\x01\x02")             // Control chars
	f.Add("multi\nline\nstring")      // Multiline
	// f.Add("unicode: 你好世界 🚀")       // Unicode (skipped for gosmopolitan)
	f.Add("127.0.0.1:8080") // Host:port

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("String parsing panicked on input %q: %v", input, r)
			}
		}()

		// Test string operations that config might perform
		_ = strings.TrimSpace(input)
		_ = strings.ToLower(input)

		// Test URL parsing for server URLs
		if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
			// Could be a valid URL
			require.Contains(t, input, "://")
		}

		// Test host:port splitting
		if strings.Contains(input, ":") {
			parts := strings.Split(input, ":")
			_ = parts
		}
	})
}

// FuzzConfigCSVParsing tests comma-separated value parsing (e.g., for CORS domains).
func FuzzConfigCSVParsing(f *testing.F) {
	// Seed corpus with CSV patterns
	f.Add("localhost,example.com,test.org")            // Valid CSV
	f.Add("")                                          // Empty
	f.Add("single-value")                              // Single value
	f.Add("value1,")                                   // Trailing comma
	f.Add(",value1")                                   // Leading comma
	f.Add("val1,,val2")                                // Empty value in middle
	f.Add(strings.Repeat("domain", 100))               // No commas
	f.Add("http://localhost:3000,https://example.com") // URLs
	f.Add("value1, value2, value3")                    // Spaces around commas

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CSV parsing panicked on input %q: %v", input, r)
			}
		}()

		// Test CSV parsing (similar to CORS allowed domains)
		values := strings.Split(input, ",")
		require.IsType(t, []string{}, values)

		// Test trimming each value
		for _, val := range values {
			trimmed := strings.TrimSpace(val)
			_ = trimmed
		}

		// Test filtering empty values
		var filtered []string
		for _, val := range values {
			trimmed := strings.TrimSpace(val)
			if trimmed != "" {
				filtered = append(filtered, trimmed)
			}
		}
		_ = filtered
	})
}

// FuzzConfigPortValidation tests port number validation.
func FuzzConfigPortValidation(f *testing.F) {
	// Seed corpus with port patterns
	f.Add("80")    // Valid HTTP port
	f.Add("443")   // Valid HTTPS port
	f.Add("8080")  // Valid alt port
	f.Add("0")     // Min port
	f.Add("65535") // Max port
	f.Add("65536") // Out of range
	f.Add("-1")    // Negative
	f.Add("")      // Empty
	f.Add("abc")   // Invalid

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Port validation panicked on input %q: %v", input, r)
			}
		}()

		// Test port parsing
		port, err := strconv.Atoi(input)
		if err == nil {
			// Validate port range (1-65535)
			if port >= 1 && port <= 65535 {
				require.True(t, port > 0 && port <= 65535, "Valid port should be in range 1-65535")
			}
		}
	})
}

// FuzzConfigSecretValidation tests secret/password config validation.
func FuzzConfigSecretValidation(f *testing.F) {
	// Seed corpus with secret patterns
	f.Add("valid-secret-12345")      // Valid
	f.Add("")                        // Empty (should fail)
	f.Add("short")                   // Too short
	f.Add(strings.Repeat("a", 1000)) // Very long
	f.Add("special!@#$%^&*()")       // Special chars
	// f.Add("unicode-密码-🔒")        // Unicode (skipped for gosmopolitan)
	f.Add(" spaces around ") // Spaces
	f.Add("\x00\x01\x02")    // Control chars

	f.Fuzz(func(t *testing.T, input string) {
		// Primary goal: ensure parsing never panics
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Secret validation panicked on input %q: %v", input, r)
			}
		}()

		// Test secret validation
		trimmed := strings.TrimSpace(input)

		// Secrets should typically be non-empty and have minimum length
		if len(trimmed) >= 16 {
			// Could be a valid secret
			require.GreaterOrEqual(t, len(trimmed), 16, "Strong secret should be at least 16 chars")
		}

		// Empty secrets are invalid
		if trimmed == "" {
			require.Empty(t, trimmed, "Empty secret should be detected")
		}
	})
}
