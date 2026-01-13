package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/xdg-go/pbkdf2"
)

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

// Encrypt encrypts the plaintext using AES-GCM.
func Encrypt(passphrase, plaintext string) (string, error) {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	_, err := rand.Read(iv)
	if err != nil {
		return "", err
	}
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		return "", err
	}
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(data), nil
}

// Decrypt decrypts the ciphertext using AES-GCM.
func Decrypt(passphrase, ciphertext string) string {
	arr := strings.Split(ciphertext, "-")
	// Validate format: must have exactly 3 components (salt-iv-data)
	if len(arr) != 3 {
		return ""
	}
	salt, err := hex.DecodeString(arr[0])
	if err != nil {
		return ""
	}
	iv, err := hex.DecodeString(arr[1])
	if err != nil {
		return ""
	}
	// IV must be exactly 12 bytes for GCM
	if len(iv) != 12 {
		return ""
	}
	data, err := hex.DecodeString(arr[2])
	if err != nil {
		return ""
	}
	key, _ := deriveKey(passphrase, salt)
	b, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		return ""
	}
	data, err = aesgcm.Open(nil, iv, data, nil)
	if err != nil {
		return ""
	}
	return string(data)
}
