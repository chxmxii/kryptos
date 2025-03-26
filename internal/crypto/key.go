package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Key []byte

func GenerateKey(size int) (Key, error) {
	if size <= 0 {
		size = 32 // Default
	}

	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}

func SaveKey(key Key, path string) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	defaultPath := filepath.Join(home + "/.kryptos/key")
	if path == "" {
		path = defaultPath
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}

	if _, err := file.Write(key); err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}

	defer file.Close()

	log.Printf("Key saved to %s\n", path)
	return nil
}

func LoadKey(path string) (Key, error) {
	// default path
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home + "/.kryptos/key")
	}

	// read key file
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	if len(key) == 0 {
		return nil, errors.New("key file is empty")
	}

	return key, nil
}

func (k Key) KeyToString() string {
	return hex.EncodeToString(k)
}

func StringToKey(s string) (Key, error) {
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	return key, nil
}
