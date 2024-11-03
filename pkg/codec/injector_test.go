package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInjector(t *testing.T) {
	tests := []struct {
		name string
		seed uint64
		rate int
		data []byte
	}{
		{
			name: "zero rate",
			seed: 42,
			rate: 0,
			data: []byte("hello world"),
		},
		{
			name: "full rate",
			seed: 42,
			rate: 100,
			data: []byte("test data"),
		},
		{
			name: "partial rate",
			seed: 42,
			rate: 500,
			data: []byte("partial injection test"),
		},
		{
			name: "empty input",
			seed: 42,
			rate: 1000,
			data: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector, err := NewInjector(tt.seed, tt.rate)
			require.NoError(t, err, "NewInjector(%d, %d)", tt.seed, tt.rate)

			encoder := injector.NewEncoder()
			decoder := injector.NewDecoder()

			t.Run("Full array encode/decode", func(t *testing.T) {
				encoded := encoder(tt.data)
				decoded := decoder(encoded)
				require.Equal(t, tt.data, decoded, "Full array encode/decode failed")
			})

			t.Run("Chunked encode/decode", func(t *testing.T) {
				encodedChunks := make([]byte, 0)
				for i := 0; i < len(tt.data); i += 2 {
					end := i + 2
					if end > len(tt.data) {
						end = len(tt.data)
					}
					chunk := encoder(tt.data[i:end])
					encodedChunks = append(encodedChunks, chunk...)
				}

				decodedChunks := make([]byte, 0)
				for i := 0; i < len(encodedChunks); i += 3 {
					end := i + 3
					if end > len(encodedChunks) {
						end = len(encodedChunks)
					}
					chunk := decoder(encodedChunks[i:end])
					decodedChunks = append(decodedChunks, chunk...)
				}

				require.Equal(t, tt.data, decodedChunks, "Chunked encode/decode failed")
			})
		})
	}
}

func TestRate(t *testing.T) {
	tests := []struct {
		name       string
		redundancy int
		want       int
	}{
		{
			name:       "zero redundancy",
			redundancy: 0,
			want:       0,
		},
		{
			name:       "100% redundancy",
			redundancy: 100,
			want:       50, // 100*100/(100+100) = 50
		},
		{
			name:       "small redundancy",
			redundancy: 10,
			want:       9, // 100*10/(100+10) ≈ 9
		},
		{
			name:       "large redundancy",
			redundancy: 900,
			want:       90, // 100*900/(100+900) = 90
		},
		{
			name:       "maximum allowed redundancy",
			redundancy: 1000,
			want:       90, // 100*1000/(100+1000) ≈ 90
		},
		{
			name:       "tiny redundancy",
			redundancy: 1,
			want:       0, // 100*1/(100+1) ≈ 0.99 -> 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rate(tt.redundancy)
			require.Equal(t, tt.want, got, "rate(%d)", tt.redundancy)
		})
	}
}
