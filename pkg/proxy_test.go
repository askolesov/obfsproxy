package pkg

import (
	"io"
	"net"
	"testing"

	"github.com/askolesov/obfsproxy/pkg/codec"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	// Create a proxy instance
	p := &Proxy{}

	// Create pipe connections
	in1, out1 := net.Pipe()
	defer in1.Close()
	defer out1.Close()
	in2, out2 := net.Pipe()
	defer in2.Close()
	defer out2.Close()

	// Test data
	testData := []byte("Hello, World!")
	expectedData := make([]byte, len(testData))
	for i, b := range testData {
		expectedData[i] = ^b // Invert bytes
	}

	codec := codec.NewInverter()

	// Run proxy in a goroutine
	go p.proxy(in2, out1, codec.NewEncoder())

	go func() {
		// Write test data to server
		_, err := in1.Write(testData)
		require.NoError(t, err, "Failed to write to server")
		in1.Close()
	}()

	// Read from client
	result := make([]byte, len(testData))
	_, err := io.ReadFull(out2, result)
	require.NoError(t, err, "Failed to read from client")

	// Compare result with expected data
	require.Equal(t, expectedData, result, "Proxy data mismatch")
}
