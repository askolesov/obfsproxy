package codec

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXorer(t *testing.T) {
	tests := []struct {
		name  string
		key   []byte
		input []byte
		want  []byte
	}{
		{
			name:  "simple single byte key",
			key:   []byte{0x01},
			input: []byte("hello"),
			want:  []byte{0x68 ^ 0x01, 0x65 ^ 0x01, 0x6C ^ 0x01, 0x6C ^ 0x01, 0x6F ^ 0x01},
		},
		{
			name:  "multi-byte key",
			key:   []byte{0x01, 0x02, 0x03},
			input: []byte("hello"),
			want:  []byte{0x68 ^ 0x01, 0x65 ^ 0x02, 0x6C ^ 0x03, 0x6C ^ 0x01, 0x6F ^ 0x02},
		},
		{
			name:  "empty input",
			key:   []byte{0x01},
			input: []byte{},
			want:  []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xorer, err := NewXorer(tt.key)
			require.NoError(t, err, "NewXorer(%v)", tt.key)

			encoder := xorer.NewEncoder()
			decoder := xorer.NewDecoder()

			// Test Encode
			encoded := encoder(tt.input)
			require.Equal(t, tt.want, encoded, "Encode() produced unexpected result")

			// Test Decode (should return to original input)
			decoded := decoder(encoded)
			require.Equal(t, tt.input, decoded, "Decode() produced unexpected result")
		})
	}
}

func TestXorerSymmetry(t *testing.T) {
	key := []byte("test-key")
	input := []byte("Hello, World!")

	xorer, err := NewXorer(key)
	require.NoError(t, err, "NewXorer(%v)", key)

	encoder := xorer.NewEncoder()
	decoder := xorer.NewDecoder()

	encoded := encoder(input)
	decoded := decoder(encoded)

	require.Equal(t, input, decoded, "Decode() produced unexpected result")

	if bytes.Equal(encoded, input) {
		t.Error("Encoded data should be different from input")
	}
}
