package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInverter(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "empty slice",
			data: []byte{},
			want: []byte{},
		},
		{
			name: "single byte",
			data: []byte{0xFF},
			want: []byte{0x00},
		},
		{
			name: "multiple bytes",
			data: []byte{0x00, 0xFF, 0xAA},
			want: []byte{0xFF, 0x00, 0x55},
		},
	}

	inverter := NewInverter()
	encoder := inverter.NewEncoder()
	decoder := inverter.NewDecoder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Encode
			encoded := encoder(tt.data)
			require.Equal(t, tt.want, encoded, "Encode() produced unexpected result")

			// Test Decode (should return to original)
			decoded := decoder(encoded)
			require.Equal(t, tt.data, decoded, "Decode() produced unexpected result")
		})
	}
}
