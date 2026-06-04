package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	secretKey []byte
	once      sync.Once
)

func getKey() []byte {
	once.Do(func() {
		key := os.Getenv("LOGCAT_SECRET_KEY")
		if key == "" {
			key = "logcat-default-secret-key-32b"
		}
		if len(key) < 32 {
			padded := make([]byte, 32)
			copy(padded, key)
			key = string(padded)
		}
		if len(key) > 32 {
			key = key[:32]
		}
		secretKey = []byte(key)
	})
	return secretKey
}

func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("invalid encrypted data")
	}
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid encrypted data")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed")
	}
	return string(plaintext), nil
}

func IsEncrypted(value string) bool {
	if value == "" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return false
	}
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return false
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return false
	}
	return len(decoded) >= aesGCM.NonceSize()
}
