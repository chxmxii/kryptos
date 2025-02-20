package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptor(t *testing.T) {

	key, err := GenerateKey(32)
	require.NoError(t, err)

	tests := []struct {
		name      string
		data      []byte
		key       Key
		wantError bool
	}{
		{
			name:      "Valid Encryption and Decryption",
			data:      []byte("Hello, World!"),
			key:       key,
			wantError: false,
		},
		{
			name:      "Empty Data",
			data:      []byte(""),
			key:       key,
			wantError: false,
		},
		{
			name:      "Invalid Key (Empty)",
			data:      []byte("Hello, World!"),
			key:       []byte{},
			wantError: true,
		},
		{
			name:      "Invalid Key (Nil)",
			data:      []byte("Hello, World!"),
			key:       nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the data
			encrypted, err := Encrypt(tt.data, tt.key)
			if tt.wantError {
				require.Error(t, err, "Expected encryption to fail")
				return
			}
			require.NoError(t, err, "Encryption should not fail")

			// Decrypt the data
			decrypted, err := Decrypt(encrypted, tt.key)
			if tt.wantError {
				require.Error(t, err, "Expected decryption to fail")
				return
			}
			require.NoError(t, err, "Decryption should not fail")

			// Verify the decrypted data matches the original
			assert.Equal(t, tt.data, decrypted, "Decrypted data should match the original")
		})
	}
}
