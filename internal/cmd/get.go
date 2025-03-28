package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chxmxii/kryptos/internal/crypto"
	"github.com/chxmxii/kryptos/internal/redis"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Example: "kryptos get role",
	Short:   "get a value from the database",
	Args:    cobra.ExactArgs(1),
	RunE:    Get,
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Add flags to the command
	getCmd.Flags().StringP("key", "k", "", "Path to the encryption key file")
	getCmd.Flags().IntP("index", "i", 0, "Redis database index")
}

func Get(cmd *cobra.Command, args []string) error {

	k := args[0]

	keyPath := cmd.Flag("key").Value.String()
	dbIndex, _ := cmd.Flags().GetInt("index")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	decryptionKey, err := crypto.LoadKey(keyPath)
	if err != nil {
		return fmt.Errorf("failed to load decryption key: %w", err)
	}

	redisClient, err := redis.Connect(redisAddr, redisPassword, dbIndex)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisClient.Close()

	// start a new context
	ctx := context.Background()

	// Check if the key exists in the database
	exists, err := redisClient.Client.Exists(ctx, k).Result()
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("key %s does not exist", k)
	}

	// Get the value from the database
	v, err := redisClient.Client.Get(ctx, k).Result()
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}

	// Decrypt the value
	decryptedValue, err := crypto.Decrypt([]byte(v), decryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}

	// Check if the decrypted value is empty
	if len(decryptedValue) == 0 {
		return fmt.Errorf("decrypted value is empty, something went wrong")
	}

	// print the decrypted value
	fmt.Printf("$> '%s'\n", decryptedValue)

	return nil
}
