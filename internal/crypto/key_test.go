package crypto

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{"Valid Key Size", 32, false},
		{"Zero Key Size", 0, false}, // Defaults to 32
		// {"Negative Key Size", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.size)
			fmt.Println(key)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, key)
				require.Equal(t, 32, len(key)) // Ensure default size is 32
			}
		})
	}
}

func TestSaveAndLoadKey(t *testing.T) {
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "test.key")

	key, err := GenerateKey(32)
	require.NoError(t, err)

	t.Run("SaveKey", func(t *testing.T) {
		err := SaveKey(key, keyPath)
		require.NoError(t, err)

		_, err = os.Stat(keyPath)
		require.NoError(t, err)
	})

	t.Run("LoadKey", func(t *testing.T) {
		loadedKey, err := LoadKey(keyPath)
		require.NoError(t, err)
		require.Equal(t, key, loadedKey)
	})

	t.Run("LoadKey_InvalidPath", func(t *testing.T) {
		_, err := LoadKey(filepath.Join(tempDir, "nonexistent.key"))
		require.Error(t, err)
	})
}

func TestSaveKey_EmptyKey(t *testing.T) {
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "empty.key")

	err := SaveKey(nil, keyPath)
	require.Error(t, err)
}

func TestKeyToString(t *testing.T) {
	key := Key{0xDE, 0xAD, 0xBE, 0xEF}
	expected := "deadbeef"

	result := key.KeyToString()
	require.Equal(t, expected, result)
}
