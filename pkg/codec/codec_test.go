package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestCodecs(t *testing.T) {
	codecs := map[string]Codec{
		"inverter": NewInverter(),
		"xorer":    must(NewXorer([]byte("test-key"))),
		"injector": must(NewInjector(42, 500)),
	}

	tests := []struct {
		name            string
		input           []byte
		encodeChunkSize int
		decodeChunkSize int
	}{
		{
			name:            "empty input",
			input:           []byte{},
			encodeChunkSize: 1,
			decodeChunkSize: 1,
		},
		{
			name:            "single byte",
			input:           []byte{42},
			encodeChunkSize: 1,
			decodeChunkSize: 1,
		},
		{
			name:            "small array whole",
			input:           []byte{1, 2, 3, 4, 5},
			encodeChunkSize: 5,
			decodeChunkSize: 5,
		},
		{
			name:            "small array different chunks",
			input:           []byte{1, 2, 3, 4, 5},
			encodeChunkSize: 2,
			decodeChunkSize: 3,
		},
		{
			name:            "random data whole vs chunks",
			input:           generateRandomBytes(1000, 1),
			encodeChunkSize: 100, // whole array
			decodeChunkSize: 7,   // chunks
		},
		{
			name:            "random data different chunks",
			input:           generateRandomBytes(1000, 2),
			encodeChunkSize: 7,
			decodeChunkSize: 13,
		},
		{
			name:            "random data 1 chunk",
			input:           generateRandomBytes(1000, 3),
			encodeChunkSize: 1, // 1 chunk
			decodeChunkSize: 1, // 1 chunk
		},
	}

	for codecName, codec := range codecs {
		t.Run(codecName, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					// Test whole array
					encoder := codec.NewEncoder()
					decoder := codec.NewDecoder()

					encoded := encoder(tt.input)
					decoded := decoder(encoded)
					require.Equal(t, tt.input, decoded, "whole array encode/decode failed")

					// Test chunks
					chunkEncoder := codec.NewEncoder() // new instances because they are stateful
					chunkDecoder := codec.NewDecoder()

					encodedChunks := transformByChunks(chunkEncoder, tt.input, tt.encodeChunkSize)
					decodedChunks := transformByChunks(chunkDecoder, encodedChunks, tt.decodeChunkSize)
					require.Equal(t, tt.input, decodedChunks, "chunked encode/decode failed")

					// Verify that encoded data is different from input (except for empty input and injector with low rate)
					if len(tt.input) > 0 && codecName != "injector" {
						require.NotEqual(t, tt.input, encoded, "encoded data should be different from input")
					}

					// Verify that chunked encoding matches whole array encoding (except for injector)
					if codecName != "injector" {
						require.Equal(t, encoded, encodedChunks, "chunked encoding should match whole array encoding")
					}
				})
			}
		})
	}
}

func generateRandomBytes(n int, seed int64) []byte {
	rng := rand.New(rand.NewSource(uint64(seed)))
	result := make([]byte, n)
	for i := range result {
		result[i] = byte(rng.Intn(256))
	}
	return result
}

func transformByChunks(transformer Transformer, input []byte, chunkSize int) []byte {
	if len(input) == 0 {
		return transformer(input)
	}

	var result []byte
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input)
		}
		chunk := input[i:end]
		decoded := transformer(chunk)
		result = append(result, decoded...)
	}
	return result
}
