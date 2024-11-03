package codec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	tests := []struct {
		name   string
		codecs []Codec
		input  []byte
	}{
		{
			name: "single codec",
			codecs: []Codec{
				NewInverter(),
			},
			input: []byte("hello"),
		},
		{
			name: "multiple codecs",
			codecs: []Codec{
				NewInverter(),
				must(NewXorer([]byte("key"))),
				must(NewInjector(42, 100)),
			},
			input: []byte("test data"),
		},
		{
			name: "empty input",
			codecs: []Codec{
				NewInverter(),
				must(NewXorer([]byte("key"))),
			},
			input: []byte{},
		},
		{
			name: "large input",
			codecs: []Codec{
				NewInverter(),
				must(NewXorer([]byte("test-key"))),
				must(NewInjector(42, 50)),
			},
			input: generateRandomBytes(1000, 42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := NewChain(tt.codecs)
			require.NoError(t, err)

			encoder := chain.NewEncoder()
			decoder := chain.NewDecoder()

			// Test full array encode/decode
			encoded := encoder(tt.input)
			decoded := decoder(encoded)
			require.Equal(t, tt.input, decoded)

			// Test chunked encode/decode
			encodedChunks := transformByChunks(encoder, tt.input, 7)
			decodedChunks := transformByChunks(decoder, encodedChunks, 13)
			require.Equal(t, tt.input, decodedChunks)
		})
	}
}

func TestChainValidation(t *testing.T) {
	_, err := NewChain(nil)
	require.Error(t, err)

	_, err = NewChain([]Codec{})
	require.Error(t, err)

	chain, err := NewChain([]Codec{NewInverter()})
	require.NoError(t, err)
	require.NotNil(t, chain)
}
