package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/MAJIAXIT/projname/api/pkg/logger"
)

type PasswordHolder interface {
	GetEncrPassword() string
	SetEncrPassword(string)
}

func Encrypt(holder PasswordHolder, plainPassword string) error {
	// Create a new hash for the encryption key
	employeePasswordKey, ok := os.LookupEnv("PASSWORD_KEY")
	if !ok {
		return logger.WrapError(errors.New("Failed to get PASSWORD_KEY from env"))
	}

	key := sha256.Sum256([]byte(employeePasswordKey))

	// Create cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return logger.WrapError(err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return logger.WrapError(err)
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return logger.WrapError(err)
	}

	// Encrypt the password
	ciphertext := gcm.Seal(nonce, nonce, []byte(plainPassword), nil)

	// Store the encrypted password
	holder.SetEncrPassword(base64.StdEncoding.EncodeToString(ciphertext))

	return nil
}

// Decrypt decrypts the password of a password holder
func Decrypt(holder PasswordHolder) (string, error) {
	// Create a new hash for the encryption key
	employeePasswordKey, ok := os.LookupEnv("PASSWORD_KEY")
	if !ok {
		return "", logger.WrapError(errors.New("Failed to get PASSWORD_KEY from env"))
	}

	key := sha256.Sum256([]byte(employeePasswordKey))

	// Decode the base64 stored password
	ciphertext, err := base64.StdEncoding.DecodeString(holder.GetEncrPassword())
	if err != nil {
		return "", logger.WrapError(err)
	}

	// Create cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", logger.WrapError(err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", logger.WrapError(err)
	}

	// Verify the length
	if len(ciphertext) < gcm.NonceSize() {
		return "", logger.WrapError(errors.New("ciphertext too short"))
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", logger.WrapError(err)
	}

	return string(plaintext), nil
}

// Compare compares an input password with the stored encrypted password
func Compare(holder PasswordHolder, inputPassword string) (match bool, err error) {
	// Get the decrypted password
	storedPassword, err := Decrypt(holder)
	if err != nil {
		return false, logger.WrapError(err)
	}

	// Compare the passwords
	return storedPassword == inputPassword, nil
}
