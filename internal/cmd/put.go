package cmd

import (
	"fmt"
	"os"

	c "github.com/chxmxii/kryptos/internal/crypto"
	"github.com/chxmxii/kryptos/internal/redis"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Put a key-value pair",
	RunE:  Put,
}

func Put(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	k, v := args[0], args[1]

	path := c.HomeDir() + "/.kryptos/test.key"
	// encrypt key and value
	key, err := c.LoadKey(path)

	if err != nil {
		return fmt.Errorf("failed to load key: %w", err)
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return fmt.Errorf("REDIS_ADDR environment variable not set")
	}

	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		return fmt.Errorf("REDIS_PASSWORD environment variable not set")
	}

	r, err := redis.Connect(addr, pass, 0)
	if err != nil {
		return err
	}
	defer r.Close()

	err = r.Client.Set(r.Client.Context(), string(k), string(v), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to put key-value pair: %w", err)
	}
	return nil
}
