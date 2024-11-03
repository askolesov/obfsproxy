package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNone(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "empty input",
			input: []byte{},
		},
		{
			name:  "simple input",
			input: []byte("hello world"),
		},
		{
			name:  "binary data",
			input: []byte{0x00, 0xFF, 0x42, 0x13, 0x37},
		},
		{
			name:  "large input",
			input: generateRandomBytes(1000, 42),
		},
	}

	none := NewNone()
	encoder := none.NewEncoder()
	decoder := none.NewDecoder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that encoder doesn't modify data
			encoded := encoder(tt.input)
			require.Equal(t, tt.input, encoded)

			// Test that decoder doesn't modify data
			decoded := decoder(encoded)
			require.Equal(t, tt.input, decoded)

			// Test with chunks
			encodedChunks := transformByChunks(encoder, tt.input, 7)
			require.Equal(t, tt.input, encodedChunks)

			decodedChunks := transformByChunks(decoder, encodedChunks, 13)
			require.Equal(t, tt.input, decodedChunks)
		})
	}
}
