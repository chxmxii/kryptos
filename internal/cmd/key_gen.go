package cmd

import (
	"fmt"

	"github.com/chxmxii/kryptos/internal/crypto"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "key-generate",
	Short: "Generate a new encryption key",
	Run:   Generate,
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("output", "o", "", "Path to the output file")
	genCmd.Flags().IntP("size", "s", 32, "Size of the key")
}

func Generate(cmd *cobra.Command, args []string) {
	outputPath := cmd.Flag("output").Value.String()

	size, err := cmd.Flags().GetInt("size")
	if err != nil {
		fmt.Printf("Failed to retrieve size flag: %v\n", err)
		return
	}

	key, err := crypto.GenerateKey(size)
	if err != nil {
		fmt.Printf("Failed to generate key: %v\n", err)
		return
	}

	if err := crypto.SaveKey(key, outputPath); err != nil {
		fmt.Printf("Failed to save key: %v\n", err)
		return
	}

	fmt.Println("Key generated and saved successfully")
}
