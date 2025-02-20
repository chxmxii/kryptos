package cmd

import (
	"fmt"
	"os"

	"github.com/chxmxii/kryptos/internal/crypto"
	"github.com/chxmxii/kryptos/internal/redis"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Put a key-value pair into Redis",
	Args:  cobra.ExactArgs(2),
	RunE:  Put,
}

func init() {
	rootCmd.AddCommand(putCmd)

	// Add flags to the command
	putCmd.Flags().StringP("key", "k", "", "Path to the encryption key file")
	putCmd.Flags().IntP("index", "i", 0, "Redis database index")
	putCmd.MarkFlagRequired("key")
}

func Put(cmd *cobra.Command, args []string) error {

	k := args[0]
	v := args[1]

	keyPath := cmd.Flag("key").Value.String()
	dbIndex, _ := cmd.Flags().GetInt("index")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	encryptionKey, err := crypto.LoadKey(keyPath)
	if err != nil {
		return fmt.Errorf("failed to load encryption key: %w", err)
	}

	redisClient, err := redis.Connect(redisAddr, redisPassword, dbIndex)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisClient.Close()

	// Encrypt the value
	encryptedValue, err := crypto.Encrypt([]byte(v), encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	if err := redisClient.Client.Set(redisClient.Client.Context(), k, encryptedValue, 0).Err(); err != nil {
		return fmt.Errorf("failed to set key-value pair in Redis: %w", err)
	}

	fmt.Println("Operation successful")
	return nil
}
