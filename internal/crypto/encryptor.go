package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(data []byte, key Key) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to read random: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	return append(nonce, ciphertext...), nil
}

func Decrypt(ct []byte, key Key) ([]byte, error) {

	if len(ct) == 0 {
		return []byte{}, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ct) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ct[:nonceSize], ct[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// func (k *Key) Encrypt(data []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(*k)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create cipher: %w", err)
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create GCM: %w", err)
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, fmt.Errorf("failed to read random: %w", err)
// 	}

// 	ciphertext := gcm.Seal(nil, nonce, data, nil)

// 	return append(nonce, ciphertext...), nil
// }

// func (k *Key) Decrypt(data []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(*k)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create cipher: %w", err)
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create GCM: %w", err)
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(data) < nonceSize {
// 		return nil, fmt.Errorf("ciphertext too short")
// 	}

// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decrypt: %w", err)
// 	}

// 	return plaintext, nil
// }
