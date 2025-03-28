package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chxmxii/kryptos/internal/redis"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "del",
	Example: "kryptos del chxmxii",
	Short:   "delete a value from the database",
	Args:    cobra.ExactArgs(1),
	RunE:    Delete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Add flags to the command
	deleteCmd.Flags().StringP("key", "k", "", "Path to the encryption key file")
	deleteCmd.Flags().IntP("index", "i", 0, "Redis database index")
}

func Delete(cmd *cobra.Command, args []string) error {

	k := args[0]

	dbIndex, _ := cmd.Flags().GetInt("index")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient, err := redis.Connect(redisAddr, redisPassword, dbIndex)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisClient.Close()

	// start a new context
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Check if the key exists in the database
	exists, err := redisClient.Client.Exists(ctx, k).Result()
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("key %s does not exist", k)
	}

	// delete the key from the database
	if err := redisClient.Client.Del(ctx, k).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", k, err)
	}

	fmt.Printf("Key %s deleted successfully\n", k)

	return nil
}
