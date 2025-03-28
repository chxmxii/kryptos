package cmd

import (
	"fmt"
	"os"

	"github.com/chxmxii/kryptos/internal/redis"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Example: "kryptos ls",
	Short:   "List all keys in the database",
	Args:    cobra.NoArgs,
	RunE:    List,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags to the command
	listCmd.Flags().IntP("index", "i", 0, "Redis database index")
}

func List(cmd *cobra.Command, args []string) error {

	dbIndex, _ := cmd.Flags().GetInt("index")
	if dbIndex < 0 {
		return fmt.Errorf("invalid Redis database index: %d", dbIndex)
	}

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
		ctx = cmd.Context()
	}

	// List all keys in the database
	keys, err := redisClient.Client.Keys(ctx, "*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		cmd.Println(key)
	}

	return nil
}
