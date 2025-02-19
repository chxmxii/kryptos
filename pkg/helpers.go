package helpers

import (
	"os"

	"fmt"
)

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	return home
}
