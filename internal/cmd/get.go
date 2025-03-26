package cmd

import (
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

	encryptionKey, err := crypto.LoadKey(keyPath)
	if err != nil {
		return fmt.Errorf("failed to load encryption key: %w", err)
	}

	redisClient, err := redis.Connect(redisAddr, redisPassword, dbIndex)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisClient.Close()
	ctx := cmd.
	// Get the value from the database
	if err := redisClient.Client.Get(k); err != nil {
		return fmt.Errorf("failed to get value from database: %w", err)
	}

	// Decrypt the value
	decryptedValue, err := crypto.Decrypt(encryptionKey, v)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}

	fmt.Println(decryptedValue)

	return nil
}
