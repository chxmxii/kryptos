package crypto

import (
	"crypto/rand"
	"fmt"
	"os"

	h "github.com/chxmxii/kryptos/pkg/helpers"
)

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

func SaveKey(key []byte) error {

	path := os.MkdirAll(os.File(h.HomeDir()+"/.kryptos/key.aes"), 0700)

	err = os.WriteFile(path, key, 0600)
	if err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}
	return nil
}
