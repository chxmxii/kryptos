package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
)

type SecretManager struct {
	client *redis.Client
	key    []byte
}

func NewSecretManager(key []byte) *SecretManager {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis password, if any
		DB:       0,                // Redis database index
	})

	// Test the connection
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic(err)
	}

	return &SecretManager{client: client, key: key}
}

func (sm *SecretManager) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	ciphertext = append(nonce, ciphertext...)
	return ciphertext, nil
}

func (sm *SecretManager) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (sm *SecretManager) SetSecret(key string, value string) error {
	encryptedValue, err := sm.encrypt([]byte(value))
	if err != nil {
		return err
	}

	err = sm.client.Set(sm.client.Context(), key, encryptedValue, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (sm *SecretManager) GetSecret(key string) (string, error) {
	encryptedValue, err := sm.client.Get(sm.client.Context(), key).Bytes()
	if err != nil {
		return "", err
	}

	decryptedValue, err := sm.decrypt(encryptedValue)
	if err != nil {
		return "", err
	}

	return string(decryptedValue), nil
}

func (sm *SecretManager) ListSecrets() ([]string, error) {
	keys, err := sm.client.Keys(sm.client.Context(), "*").Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (sm *SecretManager) DeleteSecret(key string) any {
	exists, err := sm.client.Exists(sm.client.Context(), key).Result()
	if err != nil {
		panic(err)
	}
	if exists == 1 {
		_, err := sm.client.Del(sm.client.Context(), key).Result()
		if err != nil {
			return err
		}
		fmt.Println("Secret deleted successfully!")
	} else {
		fmt.Println("Secret does not exist.")
	}
	return nil
}
func mustNot(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// specify the encryption key
	KeyPath := os.Getenv("PWD") + "/key/encryptionKey"
	encryptionKey, err := ioutil.ReadFile(KeyPath)
	mustNot(err)
	// Create a new instance of the secret manager
	secretManager := NewSecretManager(encryptionKey)

	// Define command-line flags
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Set command flags
	setKey := setCmd.String("k", "", "Key of the secret")
	setValue := setCmd.String("v", "", "Value of the secret")

	// Get and delete command flags
	getKey := getCmd.String("k", "", "Key of the secret to get")
	DeleteKey := deleteCmd.String("k", "", "Key of the secret to delete")

	// Check which subcommand is invoked
	if len(os.Args) < 2 {
		fmt.Println("Please specify a subcommand:\n \n set - Set a secret\n get - Get a secret\n list - List all the secrets\n delete - Delete a secret")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "set":
		setCmd.Parse(os.Args[2:])
		if *setKey == "" || *setValue == "" {
			fmt.Println("Please provide a key and value for the secret")
			setCmd.PrintDefaults()
			os.Exit(1)
		}
		err := secretManager.SetSecret(*setKey, *setValue)
		if err != nil {
			fmt.Println("Error setting secret:", err)
			os.Exit(1)
		}
		fmt.Println("Secret has been set successfully!")

	case "get":
		getCmd.Parse(os.Args[2:])
		if *getKey == "" {
			fmt.Println("Please provide a key for the secret")
			getCmd.PrintDefaults()
			os.Exit(1)
		}
		secret, err := secretManager.GetSecret(*getKey)
		if err != nil {
			fmt.Println("Error getting secret:", err)
			os.Exit(1)
		}
		fmt.Println("Secret:", secret)

	case "list":
		listCmd.Parse(os.Args[2:])
		secrets, err := secretManager.ListSecrets()
		if err != nil {
			fmt.Println("Error listing secrets:", err)
			os.Exit(1)
		}
		fmt.Println("My Secrets:")
		for indx, key := range secrets {
			fmt.Println(indx, key)
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *DeleteKey == "" {
			fmt.Println("Please provide a key for the secret to delete")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		err := secretManager.DeleteSecret(*DeleteKey)
		if err != nil {
			fmt.Println("Error deleting secret:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Invalid subcommand. \nUsage:\n - set\n - get\n - list\n - delete")
		os.Exit(1)
	}
}
